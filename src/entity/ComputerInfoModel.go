package entity

import (
	"strconv"
)

// ComputerInfoModel is the model of computer info
type ComputerInfoModel struct {
	CpuInfo         []CpuInfoModel  `json:"cpu_info"`
	MemInfo         MemInfoModel    `json:"mem_info"`
	HostInfo        HostInfoModel   `json:"host_info"`
	DiskInfo        []DiskInfoModel `json:"disk_info"`
	NetInfo         []NetInfoModel  `json:"net_info"`
	PoolCoreNum     int             `json:"pool_core_num"`
	PoolMaxNum      int             `json:"pool_max_num"`
	PoolActivateNum int             `json:"pool_activate_num"`
	PoolJobNum      int             `json:"pool_job_num"`
}

// CpuInfoModel is the model of cpu info
type CpuInfoModel struct {
	CpuName       string `json:"cpu_name"`
	CpuCoreNum    string `json:"cpu_core_num"`
	CpuFrequency  string `json:"cpu_frequency"`
	CpuVendorId   string `json:"cpu_vendor_id"`
	CpuPhysicalId string `json:"cpu_physical_id"`
	CpuUsage      string `json:"cpu_usage"`
}

// MemInfoModel is the model of memory info
type MemInfoModel struct {
	MemTotal string `json:"mem_total"`
	MemUsed  string `json:"mem_used"`
	MemAvi   string `json:"mem_avi"`
	MemUsage string `json:"mem_usage"`
}

// HostInfoModel is the model of host info
type HostInfoModel struct {
	HostName          string `json:"host_name"`
	HostOs            string `json:"host_os"`
	HostOsVersion     string `json:"host_os_version"`
	HostKernelVersion string `json:"host_kernel_version"`
	HostKernelArch    string `json:"host_kernel_arch"`
	HostId            string `json:"host_id"`
}

// DiskInfoModel is the model of disk info
type DiskInfoModel struct {
	DiskName  string `json:"disk_name"`
	DiskSize  string `json:"disk_size"`
	DiskUsed  string `json:"disk_used"`
	DiskAvi   string `json:"disk_avi"`
	DiskUsage string `json:"disk_usage"`
	FsType    string `json:"fs_type"`
	Opts      string `json:"opts"`
}

// NetInfoModel is the model of network info
type NetInfoModel struct {
	NetName         string `json:"net_name"`
	BytesSent       string `json:"bytes_sent"`
	BytesReceived   string `json:"bytes_received"`
	PacketsSent     string `json:"packets_sent"`
	PacketsReceived string `json:"packets_received"`
}

// NewComputerInfoModel
// @Description: create a new ComputerInfoModel
// @return       *ComputerInfoModel
func NewComputerInfoModel() *ComputerInfoModel {
	return &ComputerInfoModel{}
}

// SetComputerInfoModel
// @Description: set the ComputerInfoModel
// @receiver     com         *ComputerInfoModel
// @param        c           CpuInfoModel[]
// @param        m           MemInfoModel
// @param        h           HostInfoModel
// @param        d           DiskInfoModel[]
// @param        n           NetInfoModel[]
// @param        activateNum int
// @param        jobNum      int
func (com *ComputerInfoModel) SetComputerInfoModel(c []CpuInfoModel, m MemInfoModel,
	h HostInfoModel, d []DiskInfoModel, n []NetInfoModel, coreNum int, maxNum int, activateNum int, jobNum int) {
	com.CpuInfo = c
	com.MemInfo = m
	com.HostInfo = h
	com.DiskInfo = d
	com.NetInfo = n
	com.PoolCoreNum = coreNum
	com.PoolMaxNum = maxNum
	com.PoolActivateNum = activateNum
	com.PoolJobNum = jobNum
}

// NewCpuInfoModel
// @Description: create a new CpuInfoModel
// @param        cpuName    string
// @param        coreNum    int
// @param        freq       float64
// @param        VendorId   string
// @param        physicalId string
// @param        usage      float64
// @return       *CpuInfoModel cpu info model
func NewCpuInfoModel(cpuName string, coreNum int32, freq float64, VendorId string, physicalId string, usage float64) *CpuInfoModel {
	frequency := strconv.FormatFloat(freq/1000, 'f', 4, 64)
	coreNumber := strconv.Itoa(int(coreNum))
	usagePercent := strconv.FormatFloat(usage, 'f', 4, 64)
	return &CpuInfoModel{
		CpuName:       cpuName,
		CpuCoreNum:    coreNumber,
		CpuFrequency:  frequency,
		CpuVendorId:   VendorId,
		CpuPhysicalId: physicalId,
		CpuUsage:      usagePercent,
	}
}

// NewMemInfoModel
// @Description: create a new MemInfoModel
// @param        memTotal uint64
// @param        memAvi   uint64
// @param        memUsed  uint64
// @param        memUsage float64
// @return       *MemInfoModel memory info model
func NewMemInfoModel(memTotal uint64, memAvi uint64, memUsed uint64, memUsage float64) *MemInfoModel {
	MemTotal := strconv.FormatFloat(float64(memTotal)/1024/1024/1024, 'f', 4, 64)
	MemAvi := strconv.FormatFloat(float64(memAvi)/1024/1024/1024, 'f', 4, 64)
	MemUsed := strconv.FormatFloat(float64(memUsed)/1024/1024/1024, 'f', 4, 64)
	MemUsage := strconv.FormatFloat(memUsage, 'f', 4, 64)
	return &MemInfoModel{
		MemTotal: MemTotal,
		MemAvi:   MemAvi,
		MemUsed:  MemUsed,
		MemUsage: MemUsage,
	}
}

// NewHostInfoModel
// @Description: create a new HostInfoModel
// @param        hostName      string
// @param        os            string
// @param        osVersion     string
// @param        kernelVersion string
// @param        kernelArch    string
// @param        hostId        string
// @return       *HostInfoModel host info model
func NewHostInfoModel(hostName string, os string, osVersion string, kernelVersion string, kernelArch string, hostId string) *HostInfoModel {
	return &HostInfoModel{
		HostName:          hostName,
		HostOs:            os,
		HostOsVersion:     osVersion,
		HostKernelVersion: kernelVersion,
		HostKernelArch:    kernelArch,
		HostId:            hostId,
	}
}

// NewDiskInfoModel
// @Description: create a new DiskInfoModel
// @param        diskName  string
// @param        diskSize  uint64
// @param        diskUsed  uint64
// @param        diskAvi   uint64
// @param        diskUsage float64
// @param        fsType    string
// @param        opts      string
// @return       *DiskInfoModel disk info model
func NewDiskInfoModel(diskName string, diskSize uint64, diskUsed uint64, diskAvi uint64, diskUsage float64, fsType string, opts string) *DiskInfoModel {
	DiskSize := strconv.FormatFloat(float64(diskSize)/1024/1024/1024, 'f', 4, 64)
	DiskUsed := strconv.FormatFloat(float64(diskUsed)/1024/1024/1024, 'f', 4, 64)
	DiskAvi := strconv.FormatFloat(float64(diskAvi)/1024/1024/1024, 'f', 4, 64)
	DiskUsage := strconv.FormatFloat(diskUsage, 'f', 4, 64)
	return &DiskInfoModel{
		DiskName:  diskName,
		DiskSize:  DiskSize,
		DiskUsed:  DiskUsed,
		DiskAvi:   DiskAvi,
		DiskUsage: DiskUsage,
		FsType:    fsType,
		Opts:      opts,
	}
}

// NewNetInfoModel
// @Description: create a new NetInfoModel
// @param        name            string
// @param        bytesSent       uint64
// @param        bytesReceived   uint64
// @param        packetsSent     uint64
// @param        packetsReceived uint64
// @return       *NetInfoModel   network info model
func NewNetInfoModel(name string, bytesSent uint64, bytesReceived uint64, packetsSent uint64, packetsReceived uint64) *NetInfoModel {
	BytesSent := strconv.FormatFloat(float64(bytesSent)/1024/1024, 'f', 4, 64)
	BytesReceived := strconv.FormatFloat(float64(bytesReceived)/1024/1024, 'f', 4, 64)
	PacketsSent := strconv.FormatFloat(float64(packetsSent)/1024/1024, 'f', 4, 64)
	PacketsReceived := strconv.FormatFloat(float64(packetsReceived)/1024/1024, 'f', 4, 64)
	return &NetInfoModel{
		NetName:         name,
		BytesSent:       BytesSent,
		BytesReceived:   BytesReceived,
		PacketsSent:     PacketsSent,
		PacketsReceived: PacketsReceived,
	}
}
