package data

func NewDevice() *Device {
	d := new(Device)
	d.Dot11 = Dot11{}
	d.Dot11.AssociatedClientMap = make(map[string]string)
	return d
}

type SSID struct {
	SSID                         string  `json:"dot11.advertisedssid.ssid"`
	Len                          int     `json:"dot11.advertisedssid.ssidlen"`
	Hash                         int     `json:"dot11.advertisedssid.ssid_hash"`
	Beacon                       int     `json:"dot11.advertisedssid.beacon"`
	ProbeResponse                int     `json:"dot11.advertisedssid.probe_response"`
	Channel                      string  `json:"dot11.advertisedssid.channel"`
	HtMode                       string  `json:"dot11.advertisedssid.ht_mode"`
	HtCenter1                    int     `json:"dot11.advertisedssid.ht_center_1"`
	HtCenter2                    int     `json:"dot11.advertisedssid.ht_center_2"`
	FirstTime                    int     `json:"dot11.advertisedssid.first_time"`
	LastTime                     int     `json:"dot11.advertisedssid.last_time"`
	Cloaked                      int     `json:"dot11.advertisedssid.cloaked"`
	CryptBitfield                int64   `json:"dot11.advertisedssid.crypt_bitfield"`
	CryptSet                     int     `json:"dot11.advertisedssid.crypt_set"`
	CryptString                  string  `json:"dot11.advertisedssid.crypt_string"`
	Maxrate                      float64 `json:"dot11.advertisedssid.maxrate"`
	Beaconrate                   int     `json:"dot11.advertisedssid.beaconrate"`
	BeaconsSec                   int     `json:"dot11.advertisedssid.beacons_sec"`
	IetagChecksum                int     `json:"dot11.advertisedssid.ietag_checksum"`
	WpaMfpRequired               int     `json:"dot11.advertisedssid.wpa_mfp_required"`
	WpaMfpSupported              int     `json:"dot11.advertisedssid.wpa_mfp_supported"`
	Dot11RMobility               int     `json:"dot11.advertisedssid.dot11r_mobility"`
	Dot11RMobilityDomainId       int     `json:"dot11.advertisedssid.dot11r_mobility_domain_id"`
	Dot11EQbss                   int     `json:"dot11.advertisedssid.dot11e_qbss"`
	Dot11EQbssStations           int     `json:"dot11.advertisedssid.dot11e_qbss_stations"`
	Dot11EChannelUtilizationPerc float64 `json:"dot11.advertisedssid.dot11e_channel_utilization_perc"`
	CcxTxpower                   int     `json:"dot11.advertisedssid.ccx_txpower"`
	CiscoClientMfp               int     `json:"dot11.advertisedssid.cisco_client_mfp"`
	AdvertisedTxpower            int     `json:"dot11.advertisedssid.advertised_txpower"`
	Dot11DCountry                string  `json:"dot11.advertisedssid.dot11d_country"`
}

type Dot11 struct {
	Typeset                int    `json:"dot11.device.typeset"`
	NumClientAps           int    `json:"dot11.device.num_client_aps"`
	NumAdvertisedSsids     int    `json:"dot11.device.num_advertised_ssids"`
	NumRespondedSsids      int    `json:"dot11.device.num_responded_ssids"`
	NumProbedSsids         int    `json:"dot11.device.num_probed_ssids"`
	NumAssociatedClients   int    `json:"dot11.device.num_associated_clients"`
	ClientDisconnects      int    `json:"dot11.device.client_disconnects"`
	ClientDisconnectsLast  int    `json:"dot11.device.client_disconnects_last"`
	LastSequence           int    `json:"dot11.device.last_sequence"`
	BssTimestamp           int64  `json:"dot11.device.bss_timestamp"`
	NumFragments           int    `json:"dot11.device.num_fragments"`
	NumRetries             int    `json:"dot11.device.num_retries"`
	Datasize               int    `json:"dot11.device.datasize"`
	DatasizeRetry          int    `json:"dot11.device.datasize_retry"`
	LastBeaconTimestamp    int    `json:"dot11.device.last_beacon_timestamp"`
	WpsM3Count             int    `json:"dot11.device.wps_m3_count"`
	WpsM3Last              int    `json:"dot11.device.wps_m3_last"`
	MinTxPower             int    `json:"dot11.device.min_tx_power"`
	MaxTxPower             int    `json:"dot11.device.max_tx_power"`
	LinkMeasurementCapable int    `json:"dot11.device.link_measurement_capable"`
	NeighborReportCapable  int    `json:"dot11.device.neighbor_report_capable"`
	BeaconFingerprint      int    `json:"dot11.device.beacon_fingerprint"`
	ProbeFingerprint       int    `json:"dot11.device.probe_fingerprint"`
	ResponseFingerprint    int    `json:"dot11.device.response_fingerprint"`
	LastBssid              string `json:"dot11.device.last_bssid"`
	/*
		"dot11.device.associated_client_map": {
		      	"00:00:00:E2:00:35": "0700000D00270204_3E5AE00B1F72",
		      	"00:00:00:7B:00:18": "D240020000007007_B7DC1C1E768B",
		      	"00:00:00:44:00:93": "20000740D0007020_547C934F4693",
		},
	*/
	AssociatedClientMap    map[string]string `json:"dot11.device.associated_client_map"`
	AdvertisedSsidMap      []SSID            `json:"dot11.device.advertised_ssid_map"`
	LastBeaconedSsidRecord SSID              `json:"dot11.device.last_beaconed_ssid_record"`
}

type Device struct {
	Dot11               Dot11  `json:"dot11.device"`
	BaseKey             string `json:"kismet.device.base.key"`
	BaseMacaddr         string `json:"kismet.device.base.macaddr"`
	BaseRelatedDevices  any    // TODO: wat schema dis, Dot11?
	BaseName            string `json:"kismet.device.base.name"`
	BaseCommonname      string `json:"kismet.device.base.commonname"`
	ServerUuid          string `json:"kismet.server.uuid"`
	BaseBasicTypeSet    int    `json:"kismet.device.base.basic_type_set"`
	BaseCrypt           string `json:"kismet.device.base.crypt"`
	BaseBasicCryptSet   int    `json:"kismet.device.base.basic_crypt_set"`
	BaseFirstTime       int    `json:"kismet.device.base.first_time"`
	BaseLastTime        int    `json:"kismet.device.base.last_time"`
	BaseModTime         int    `json:"kismet.device.base.mod_time"`
	BasePacketsTotal    int    `json:"kismet.device.base.packets.total"`
	BasePacketsRxTotal  int    `json:"kismet.device.base.packets.rx_total"`
	BasePacketsTxTotal  int    `json:"kismet.device.base.packets.tx_total"`
	BasePacketsLlc      int    `json:"kismet.device.base.packets.llc"`
	BasePacketsError    int    `json:"kismet.device.base.packets.error"`
	BasePacketsData     int    `json:"kismet.device.base.packets.data"`
	BasePacketsCrypt    int    `json:"kismet.device.base.packets.crypt"`
	BasePacketsFiltered int    `json:"kismet.device.base.packets.filtered"`
	BaseDatasize        int    `json:"kismet.device.base.datasize"`
	BaseChannel         string `json:"kismet.device.base.channel"`
	BaseFrequency       int    `json:"kismet.device.base.frequency"`
	BaseFreqKhzMap      any    // TODO

	BaseNumAlerts int `json:"kismet.device.base.num_alerts"`
	BaseSeenby    []struct {
		CommonSeenbyFirstTime  int    `json:"kismet.common.seenby.first_time"`
		CommonSeenbyLastTime   int    `json:"kismet.common.seenby.last_time"`
		CommonSeenbyNumPackets int    `json:"kismet.common.seenby.num_packets"`
		CommonSeenbyUuid       string `json:"kismet.common.seenby.uuid"`
	} `json:"kismet.device.base.seenby"`
	BasePhyname string `json:"kismet.device.base.phyname"`
	BaseManuf   string `json:"kismet.device.base.manuf"`
	BaseSignal  struct {
		CommonSignalType        string `json:"kismet.common.signal.type"`
		CommonSignalLastSignal  int    `json:"kismet.common.signal.last_signal"`
		CommonSignalLastNoise   int    `json:"kismet.common.signal.last_noise"`
		CommonSignalMinSignal   int    `json:"kismet.common.signal.min_signal"`
		CommonSignalMinNoise    int    `json:"kismet.common.signal.min_noise"`
		CommonSignalMaxSignal   int    `json:"kismet.common.signal.max_signal"`
		CommonSignalMaxNoise    int    `json:"kismet.common.signal.max_noise"`
		CommonSignalMaxseenrate int    `json:"kismet.common.signal.maxseenrate"`
		CommonSignalEncodingset int    `json:"kismet.common.signal.encodingset"`
		CommonSignalCarrierset  int    `json:"kismet.common.signal.carrierset"`
		CommonSignalSignalRrd   struct {
			CommonRrdLastTime    int   `json:"kismet.common.rrd.last_time"`
			CommonRrdSerialTime  int   `json:"kismet.common.rrd.serial_time"`
			CommonRrdLastValue   int   `json:"kismet.common.rrd.last_value"`
			CommonRrdLastValueN1 int   `json:"kismet.common.rrd.last_value_n1"`
			CommonRrdMinuteVec   []int `json:"kismet.common.rrd.minute_vec"`
			CommonRrdBlankVal    int   `json:"kismet.common.rrd.blank_val"`
		} `json:"kismet.common.signal.signal_rrd"`
	} `json:"kismet.device.base.signal"`
	BaseTxPacketsRrd struct {
		CommonRrdLastTime    int       `json:"kismet.common.rrd.last_time"`
		CommonRrdSerialTime  int       `json:"kismet.common.rrd.serial_time"`
		CommonRrdLastValue   int       `json:"kismet.common.rrd.last_value"`
		CommonRrdLastValueN1 int       `json:"kismet.common.rrd.last_value_n1"`
		CommonRrdMinuteVec   []int     `json:"kismet.common.rrd.minute_vec"`
		CommonRrdHourVec     []float64 `json:"kismet.common.rrd.hour_vec"`
		CommonRrdDayVec      []float64 `json:"kismet.common.rrd.day_vec"`
		CommonRrdBlankVal    int       `json:"kismet.common.rrd.blank_val"`
	} `json:"kismet.device.base.tx_packets.rrd"`
	BasePacketsRrd struct {
		CommonRrdLastTime    int       `json:"kismet.common.rrd.last_time"`
		CommonRrdSerialTime  int       `json:"kismet.common.rrd.serial_time"`
		CommonRrdLastValue   int       `json:"kismet.common.rrd.last_value"`
		CommonRrdLastValueN1 int       `json:"kismet.common.rrd.last_value_n1"`
		CommonRrdMinuteVec   []int     `json:"kismet.common.rrd.minute_vec"`
		CommonRrdHourVec     []float64 `json:"kismet.common.rrd.hour_vec"`
		CommonRrdDayVec      []float64 `json:"kismet.common.rrd.day_vec"`
		CommonRrdBlankVal    int       `json:"kismet.common.rrd.blank_val"`
	} `json:"kismet.device.base.packets.rrd"`
	BaseDatasizeRrd struct {
		CommonRrdLastTime    int       `json:"kismet.common.rrd.last_time"`
		CommonRrdSerialTime  int       `json:"kismet.common.rrd.serial_time"`
		CommonRrdLastValue   int       `json:"kismet.common.rrd.last_value"`
		CommonRrdLastValueN1 int       `json:"kismet.common.rrd.last_value_n1"`
		CommonRrdMinuteVec   []int     `json:"kismet.common.rrd.minute_vec"`
		CommonRrdHourVec     []float64 `json:"kismet.common.rrd.hour_vec"`
		CommonRrdDayVec      []float64 `json:"kismet.common.rrd.day_vec"`
		CommonRrdBlankVal    int       `json:"kismet.common.rrd.blank_val"`
	} `json:"kismet.device.base.datasize.rrd"`
	BaseType string `json:"kismet.device.base.type"`
}
