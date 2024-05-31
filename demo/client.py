import yaml
import logging
import time
import requests
import socket
import queue
import threading
from typing import TypedDict
import json
import select
import sys

SUCCESS_CODE = 20010
FAILURE_CODE = 20011

class Server:
    ip: str
    port: int
    name: str
    colony: str
    namespace: str
    udp_port: int

# enum Leader and Follower
class Role:
    Leader = 1
    Follower = 2


class Config:
    def __init__(self):
        """
        Load the config file and set the client config
        """
        try:
            with open('setup.yaml', encoding='utf-8') as f:
                config = yaml.load(f, Loader=yaml.FullLoader)

                # client config
                self.ip = config['client']['ip']
                self.port = config['client']['port']
                self.name = config['client']['name']
                self.colony = config['client']['colony']
                self.namespace = config['client']['namespace']
                self.udp_port = config['client']['udp-listen']
                self.enable_blockchain = config['client']['enable-blockchain']

                # themis config
                self.server_ip = config['Themis']['ip']
                self.server_port = config['Themis']['port']
                self.beat_enable = config['Themis']['beat-enable']
                self.beat_timeout = int(config['Themis']['beat-timeout'])
                self.election_timeout = int(config['Themis']['election-timeout'])

        except FileNotFoundError:
            print('Config file not found!')
            exit(1)

        self.log = logging.getLogger('client')
        self.log.setLevel(logging.DEBUG)
        self.log.addHandler(logging.StreamHandler())
        # set format for log
        formatter = logging.Formatter('%(asctime)s | [%(name)s] | [%(levelname)s] => %(message)s')
        self.log.handlers[0].setFormatter(formatter)

        self.log.info('Client config loaded successfully! Details:')
        self.log.info('  - IP: %s', self.ip)
        self.log.info('  - Port: %s', self.port)
        self.log.info('  - Name: %s', self.name)
        self.log.info('  - Colony: %s', self.colony)
        self.log.info('  - Namespace: %s', self.namespace)
        self.log.info('  - UDP Port: %s', self.udp_port)
        self.log.info('  - Enable Blockchain: %s', self.enable_blockchain)
        self.log.info('Themis config loaded successfully! Details:')
        self.log.info('  - Themis IP: %s', self.server_ip)
        self.log.info('  - Themis Port: %s', self.server_port)
        self.log.info('  - Beat Enable: %s', self.beat_enable)
        self.log.info('  - Beat Timeout: %s', self.beat_timeout)
        self.log.info('  - Election Timeout: %s', self.election_timeout)

        self.udp_socket = None

        if self.__create_udp_listen():
            # create multi-thread to listen UDP
            self.channel = queue.Queue()
            self.udp_thread = threading.Thread(target=self.__udp_listen)
            self.udp_thread.start()

            self.log.info('Client started successfully!')
        else:
            self.log.error('Client started failed!')
            exit(1)

    def __create_udp_listen(self) -> bool:
        """
        Create a UDP listen
        """
        try:
            self.udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            # Allow reusing the same address
            self.udp_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            self.udp_socket.bind((self.ip, self.udp_port))
            self.log.info('UDP listen created successfully! -> %s:%s', self.ip, self.udp_port)
            return True
        except Exception as e:
            self.log.error('UDP listen created failed! Error: %s', e)
            return False
        return False
    
    def __udp_listen(self):
        """
        Listen UDP message
        """
        while True:
            try:
                data, addr = self.udp_socket.recvfrom(1024)
                self.log.debug('Received message from %s:%s -> %s', addr[0], addr[1], data.decode('utf-8'))
                server = json.loads(data.decode('utf-8'))
                self.channel.put({
                    "ip": server['IP'],
                    "port": int(server['port']),
                    "name": server['name'],
                    "colony": server['colony'],
                    "namespace": server['namespace'],
                    "udp_port": server['udp_port']
                })
            except Exception as e:
                self.log.error('UDP listen failed! Error: %s', e)
                break

    def __del__(self):
        """
        Close the UDP listen
        """
        if self.udp_socket:
            self.udp_socket.close()
            self.log.info('UDP listen closed successfully!')

    def registe(self) -> bool:
        """
        POST http://124.220.162.209:7480/v1/message/leader/register
        {
            "IP": "0.0.0.0",
            "colony": "",
            "name": "默认",
            "namespace": "",
            "port": "80",
            "time": "2024.05.31",
            "udp_port": "50000"
        }
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/leader/register'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Registering... -> %s', data)
            response = requests.post(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Register success! -> %s', response.json())
                else:
                    self.log.error('Register failed! -> %s', response.json())
                return True
            else:
                self.log.error('Register failed! -> %s', response.json())
                return False
        except Exception as e:
            self.log.error('Register failed! Error: %s', e)
            return False

    def beat(self) -> bool:
        """
        PUT http://124.220.162.209:7480/v1/message/leader/beat
                {
            "IP": "0.0.0.0",
            "colony": "",
            "name": "默认",
            "namespace": "",
            "port": "80",
            "time": "2024.05.31",
            "udp_port": "50000"
        }
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/leader/beat'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Beating... -> %s', data)
            response = requests.put(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Beat success! -> %s', response.json())
                else:
                    self.log.error('Beat failed! -> %s', response.json())
                return True
            else:
                self.log.error('Beat failed! -> %s', response.json())
                return False
        except Exception as e:
            self.log.error('Beat failed! Error: %s', e)
            return False

    def getLeader(self) -> Server:
        """
        POST http://124.220.162.209:7480/v1/message/follow/getLeader
        {
            "IP": "192.168.2.1",
            "colony": "西安",
            "name": "西安-服务2",
            "namespace": "西北dev",
            "port": "80",
            "time": "2024.05.31",
            "udp_port": "50000"
        }
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/follow/getLeader'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Getting leader... -> %s', data)
            response = requests.post(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Get leader success! -> %s', response.json())
                    leader = response.json()['data']
                    return {
                        'ip': leader['IP'],
                        'port': leader['port'],
                        'name': leader['name'],
                        'colony': leader['colony'],
                        'namespace': leader['namespace'],
                        'udp_port': leader['udp_port']
                    }
                else:
                    self.log.error('Get leader failed! -> %s', response.json())
                    return None
            else:
                self.log.error('Get leader failed! -> %s', response.json())
                return None
        except Exception as e:
            self.log.error('Get leader failed! Error: %s', e)
            return None

    def getServers(self) -> []:
        """
        POST http://124.220.162.209:7480/v1/message/follow/getServers
        [
            {
                "IP": "192.168.2.1",
                "colony": "西安",
                "name": "西安-服务2",
                "namespace": "西北dev",
                "port": "80",
                "time": "2024.05.31",
                "udp_port": "50000"
            }
        ]
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/follow/getServers'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Getting servers... -> %s', data)
            response = requests.post(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Get servers success! -> %s', response.json())
                    servers = response.json()['data']
                    return [{
                        'ip': server['IP'],
                        'port': server['port'],
                        'name': server['name'],
                        'colony': server['colony'],
                        'namespace': server['namespace'],
                        'udp_port': server['udp_port']
                    } for server in servers]
                else:
                    self.log.error('Get servers failed! -> %s', response.json())
                    return []
            else:
                self.log.error('Get servers failed! -> %s', response.json())
                return []
        except Exception as e:
            self.log.error('Get servers failed! Error: %s', e)
            return []

    def getServersNum(self) -> int:
        """
        POST http://124.220.162.209:7480/v1/message/follow/getServersNum
        {
            "IP": "192.168.2.1",
            "colony": "西安",
            "name": "西安-服务2",
            "namespace": "西北dev",
            "port": "80",
            "time": "2024.05.31",
            "udp_port": "50000"
        }
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/follow/getServersNum'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Getting servers number... -> %s', data)
            response = requests.post(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Get servers number success! -> %s', response.json())
                    return response.json()['data']
                else:
                    self.log.error('Get servers number failed! -> %s', response.json())
                    return -1
            else:
                self.log.error('Get servers number failed! -> %s', response.json())
                return -1
        except Exception as e:
            self.log.error('Get servers number failed! Error: %s', e)
            return -1

    def election(self) -> bool:
        """
        PUT http://124.220.162.209:7480/v1/message/leader/election
        {
            "IP": "192.168.2.1",
            "colony": "西安",
            "name": "西安-服务2",
            "namespace": "西北dev",
            "port": "80",
            "time": "2024.05.31",
            "udp_port": "50000"
        }
        
        """
        url = 'http://' + str(self.server_ip) + ':' + str(self.server_port) + '/v1/message/leader/election'
        data = {
            'IP': self.ip,
            'colony': self.colony,
            'name': self.name,
            'namespace': self.namespace,
            'port': str(self.port),
            # 'time': time.strftime('%Y.%m.%d', time.localtime()),
            'udp_port': str(self.udp_port)
        }
        try:
            self.log.debug('Election...')
            response = requests.put(url, json=data)
            if response.status_code == 200:
                self.log.debug('Send Success! -> %s', response.json())
                if response.json()['code'] == SUCCESS_CODE:
                    self.log.info('Election success! -> %s', response.json())
                else:
                    self.log.error('Election failed! -> %s', response.json())
                return True
            else:
                self.log.error('Election failed! -> %s', response.json())
                return False
        except Exception as e:
            self.log.error('Election failed! Error: %s', e)
            return False


if __name__ == '__main__':
    config = Config()

    # registe
    if not config.registe():
        exit(1)

    time.sleep(1)

    # get now leader
    now_leader = config.getLeader()

    last_beat_time = time.time()
    last_update_leader_time = time.time()

    # flag
    flag = Role.Follower
    
    while True:
        if config.beat_enable:
            # heartbeat
            if time.time() - last_beat_time > config.beat_timeout:
                config.beat()
                last_beat_time = time.time()
        
        if config.channel.qsize() > 0:
            # get new leader
            now_leader = config.channel.get()
            last_update_leader_time = time.time()
            config.log.info('Receive new leader! -> %s', now_leader)

            if now_leader['ip'] == config.ip and now_leader['port'] == config.port:
                flag = Role.Leader
                servers = config.getServers()
                config.log.info('Get followers! -> %s', servers)
            else:
                flag = Role.Follower
            
            # get all servers
            config.log.info('-    Number of Servers: {}'.format(config.getServersNum()))
        
        if flag == Role.Follower and time.time() - last_update_leader_time > config.election_timeout:
            # follower election
            if not config.election():
                exit(1)
            last_update_leader_time = time.time()

        time.sleep(1)
        config.log.debug('Next loop... | Role as %s | %s', 'Leader' if flag == Role.Leader else 'Follower', time.strftime('%Y-%m-%d %H:%M:%S', time.localtime()))
        if flag == Role.Leader: # show all followers
            output = "     Followers -> "
            for server in servers:
                output += format('%s:%s [%s] | ', server['ip'], server['port'], server['name'], end='')
            config.log.debug(output)
        else: # show leader
            config.log.debug('     Leader -> %s:%s [%s]', now_leader['ip'], now_leader['port'], now_leader['name'])