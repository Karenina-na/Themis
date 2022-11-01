package entity

// ComputerInfoModel is the model of computer info
type ComputerInfoModel struct {
	CpuInfo         []CpuInfoModel  `json:"cpu_info"`
	MemInfo         MemInfoModel    `json:"mem_info"`
	HostInfo        HostInfoModel   `json:"host_info"`
	DiskInfo        []DiskInfoModel `json:"disk_info"`
	NetInfo         []NetInfoModel  `json:"net_info"`
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
