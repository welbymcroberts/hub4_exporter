package collectors

import (
	"github.com/tidwall/gjson"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"hub4_exporter/config"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)


type Exporter struct {
	mutex  sync.Mutex
	config *config.Config

	// Status
	scrapeStatus       *prometheus.Desc
	aquiredDSChannel   *prometheus.Desc
	rangedUSChannel    *prometheus.Desc
	provisioningStatus *prometheus.Desc
	BPIState           *prometheus.Desc
	maxCPE             *prometheus.Desc
	networkAccess      *prometheus.Desc
	DSFlowID           *prometheus.Desc
	DOCSISVersion      *prometheus.Desc
	USFlowID           *prometheus.Desc
	DSTrafficRate      *prometheus.Desc
	USTrafficRate      *prometheus.Desc
	DSTrafficRateMin   *prometheus.Desc
	USTrafficRateMin   *prometheus.Desc
	DSTrafficRateBurst *prometheus.Desc
	USTrafficRateBurst *prometheus.Desc
	USTrafficConnBurst *prometheus.Desc
	DSChannelPostRS    *prometheus.Desc
	DSChannelPreRS     *prometheus.Desc
	DSChannelSNR       *prometheus.Desc
	DSChannelLocked    *prometheus.Desc
	DSChannelPower     *prometheus.Desc
	DSChannelRXMer     *prometheus.Desc
	DSNumber31 *prometheus.Desc
	DSNumber *prometheus.Desc
	USNumber *prometheus.Desc
	USNumber31 *prometheus.Desc
	USChannelPower *prometheus.Desc
	USChannelSymbolRate* prometheus.Desc
	USChannelTimeouts *prometheus.Desc
	DS31ChannelLocked *prometheus.Desc
	DS31ChannelPLCPower *prometheus.Desc
	DS31ChannelRXMer *prometheus.Desc
	DS31ChannelPreRS *prometheus.Desc
	DS31ChannelPostRS *prometheus.Desc
	DS31ChannelFirstSubcarrier *prometheus.Desc
	DS31ChannelSubcarriers *prometheus.Desc
	DS31ChannelWidth *prometheus.Desc

}
var namespace = "hub4"

func PromExporter(timeout time.Duration, conf *config.Config) *Exporter {
	return &Exporter{
		config: conf,
		scrapeStatus: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"scrape",
				"status",
			),
			"Scrape Status",
			[]string{"instance", "address"},
			nil,
		),
		aquiredDSChannel: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"aquired_DS_channel",
			),
			"Acquired Downstream Channel Frequency (Hz)",
			[]string{"instance", "address"},
			nil,
		),
		rangedUSChannel: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ranged_US_channel",
			),
			"Ranged (discovered) Upstream Channel Frequency (Hz)",
			[]string{"instance", "address"},
			nil,
		),
		provisioningStatus: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"provisioning_status",
			),
			"Provisioning Status",
			[]string{"instance", "address"},
			nil,
		),
		networkAccess: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"network_access",
			),
			"Network Access",
			[]string{"instance", "address"},
			nil,
		),
		maxCPE: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"max_cpe",
			),
			"Maximum CPE devices",
			[]string{"instance", "address"},
			nil,
		),
		BPIState: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"bpi_state",
			),
			"BPI State",
			[]string{"instance", "address"},
			nil,
		),
		DOCSISVersion: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"docsis_version",
			),
			"Docsis version",
			[]string{"instance", "address"},
			nil,
		),
		DSFlowID: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_flow_id",
			),
			"DS Flow ID",
			[]string{"instance", "address"},
			nil,
		),
		USFlowID: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_flow_id",
			),
			"US Flow ID",
			[]string{"instance", "address"},
			nil,
		),
		DSTrafficRate: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_traffic_rate",
			),
			"DS Traffic Rate (max)",
			[]string{"instance", "address"},
			nil,
		),
		USTrafficRate: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_traffic_rate",
			),
			"US Traffic Rate (max)",
			[]string{"instance", "address"},
			nil,
		),
		DSTrafficRateBurst: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_traffic_rate_burst",
			),
			"DS Traffic Rate (burst)",
			[]string{"instance", "address"},
			nil,
		),
		USTrafficRateBurst: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_traffic_rate_burst",
			),
			"US Traffic Rate (busrt)",
			[]string{"instance", "address"},
			nil,
		),
		DSTrafficRateMin: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_traffic_rate_min",
			),
			"DS Traffic Rate (min)",
			[]string{"instance", "address"},
			nil,
		),
		USTrafficRateMin: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_traffic_rate_min",
			),
			"US Traffic Rate (min)",
			[]string{"instance", "address"},
			nil,
		),
		USTrafficConnBurst: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_concatenated_burst",
			),
			"US Concatenated Burst",
			[]string{"instance", "address"},
			nil,
		),
		DSChannelPower: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_power",
			),
			"DS Channel Power (dBmV)",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		DSChannelSNR: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_snr",
			),
			"DS Channel SNR (dB)",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		DSChannelLocked: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_locked",
			),
			"DS Channel Locked",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		DSChannelPreRS: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_prers_errors",
			),
			"DS Channel Recoverable Errors (Pre RS)",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		DSChannelPostRS: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_postrs_errors",
			),
			"DS Channel Unrecoverable Errors (Post RS)",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		DSChannelRXMer: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_rxmer",
			),
			"DS Channel RXMer (dB)",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		USNumber: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_channel_count",
			),
			"US Channel count",
			[]string{"instance", "address"},
			nil,
		),
		DSNumber: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds_channel_count",
			),
			"DS Channel count",
			[]string{"instance", "address"},
			nil,
		),
		USNumber31: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us31_channel_count",
			),
			"US 3.1 Channel count",
			[]string{"instance", "address"},
			nil,
		),
		DSNumber31: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_count",
			),
			"DS 3.1 Channel count",
			[]string{"instance", "address"},
			nil,
		),
		USChannelPower: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_channel_power",
			),
			"US Channel Power dBmV",
			[]string{"instance", "address", "frequency"},
			nil,
		),
		USChannelTimeouts: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"us_channel_timeouts",
			),
			"US Channel Timeouts",
			[]string{"instance", "address", "frequency", "timeout_class"},
			nil,
		),

		DS31ChannelRXMer: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_rxmer",
			),
			"DS 3.1 Channel RXMer (dB)",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelPLCPower: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_plc_power",
			),
			"DS 3.1 Channel PLC Power (dBmV)",
			[]string{"instance", "address", "id"},
			nil,
		),


		DS31ChannelLocked: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_locked",
			),
			"DS 3.1 Channel Locked",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelPreRS: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_prers_errors",
			),
			"DS 3.1 Channel Recoverable Errors (Pre RS)",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelPostRS: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_postrs_errors",
			),
			"DS 3.1 Channel Unrecoverable Errors (Pre RS)",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelFirstSubcarrier: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_first_subcarrier",
			),
			"DS 3.1 Channel First Subcarrier (Hz)",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelSubcarriers: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_subcarriers",
			),
			"DS 3.1 Channel Subcarriers",
			[]string{"instance", "address", "id"},
			nil,
		),
		DS31ChannelWidth: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"",
				"ds31_channel_width",
			),
			"DS 3.1 Channel Width (MHz)",
			[]string{"instance", "address", "id"},
			nil,
		),


	}
}


func (p *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.scrapeStatus
}

func (p *Exporter) Collect(ch chan<- prometheus.Metric) {

	// Lock so no more than 1 collect occurs at once
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Create a wait group the size of the number of configured instances
	instanceWG := sync.WaitGroup{}
	instanceWG.Add(len(p.config.Instances))

	for _, instance := range p.config.Instances {
		go func(instance *config.InstancesConfig) {
			log.Infof("Collecting for instance path: %s", instance.Name)
			// Make a HTTP client
			var httpClient = &http.Client{
				Timeout: time.Second * 30,
			}
			// Get Docsis Stats
			response, err := httpClient.Get(fmt.Sprintf("http://%s/php/ajaxGet_device_networkstatus_data.php", instance.Address))
			if err != nil {
				fmt.Println(err)
			}

			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			// Scrape Status
			ch <- prometheus.MustNewConstMetric(p.scrapeStatus, prometheus.GaugeValue, float64(1), instance.Name, instance.Address)
			// Data is
			// 0 - Acquired DS Channel

			value := gjson.Get(string(body), "0")
			ch <- prometheus.MustNewConstMetric(p.aquiredDSChannel, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 1 - Ranged US Channel
			value = gjson.Get(string(body), "1")
			ch <- prometheus.MustNewConstMetric(p.rangedUSChannel, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 2 - Acquired DS Channel Status
			// TODO: Map Text to Numeric?
			//value = gjson.Get(string(body), "2")
			//fmt.Printf("AQ DS Channel Status: %s\n", value)

			// 3 - Ranged US Channel status
			// TODO: Map Text to Numeric?
			//value = gjson.Get(string(body), "3")
			//fmt.Printf("Ranged DS Channel Status: %s\n", value)

			// 4 - Provisioning State
			value = gjson.Get(string(body), "4")
			ch <- prometheus.MustNewConstMetric(p.provisioningStatus, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)


			// 5 - Network Access
			value = gjson.Get(string(body), "5")
			if value.String() == "true" {
				ch <- prometheus.MustNewConstMetric(p.networkAccess, prometheus.GaugeValue, float64(1), instance.Name, instance.Address)
			} else {
				ch <- prometheus.MustNewConstMetric(p.networkAccess, prometheus.GaugeValue, float64(0), instance.Name, instance.Address)
			}


			// 6 - Max CPE Allowed
			value = gjson.Get(string(body), "6")
			ch <- prometheus.MustNewConstMetric(p.maxCPE, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 7 - BPI State
			value = gjson.Get(string(body), "7")
			if value.String() == "true" {
				ch <- prometheus.MustNewConstMetric(p.BPIState, prometheus.GaugeValue, float64(1), instance.Name, instance.Address)
			} else {
				ch <- prometheus.MustNewConstMetric(p.BPIState, prometheus.GaugeValue, float64(0), instance.Name, instance.Address)
			}

			// 8 - Docsis Version
			value = gjson.Get(string(body), "8")
			ch <- prometheus.MustNewConstMetric(p.DOCSISVersion, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 9 - Boot File
			//value = gjson.Get(string(body), "9")
			//fmt.Printf("Config File: %s\n", value)

			// 10 - DS Flow ID
			value = gjson.Get(string(body), "10")
			ch <- prometheus.MustNewConstMetric(p.DSFlowID, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 11 - DS Traffic Rate (Max)
			value = gjson.Get(string(body), "11")
			ch <- prometheus.MustNewConstMetric(p.DSTrafficRate, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 12 - DS Traffic Rate (Max burst)
			value = gjson.Get(string(body), "12")
			ch <- prometheus.MustNewConstMetric(p.DSTrafficRateBurst, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 13 - DS Min Traffic Rate
			value = gjson.Get(string(body), "13")
			ch <- prometheus.MustNewConstMetric(p.DSTrafficRateMin, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 14 - US Flow ID
			value = gjson.Get(string(body), "14")
			ch <- prometheus.MustNewConstMetric(p.USFlowID, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 15 - US Traffic Rate (Max)
			value = gjson.Get(string(body), "15")
			ch <- prometheus.MustNewConstMetric(p.USTrafficRate, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 16 - US Trafic Rate (Max burst)
			value = gjson.Get(string(body), "16")
			ch <- prometheus.MustNewConstMetric(p.USTrafficRateBurst, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 17 - US Min Trafic Rate
			value = gjson.Get(string(body), "17")
			ch <- prometheus.MustNewConstMetric(p.USTrafficRateMin, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 18 - US Max Conn Burst
			value = gjson.Get(string(body), "18")
			ch <- prometheus.MustNewConstMetric(p.USTrafficConnBurst, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 19 - Scheduling Type
			//value = gjson.Get(string(body), "19")
			//fmt.Printf("Scheduling Type: %s\n", value)

			// 20 - DS Channel - JSON
			dsChannels := gjson.Parse(gjson.Get(string(body), "20").String())
			dsChannels.ForEach(func(key, value gjson.Result) bool {
				channel := gjson.Parse(value.String())
				//id := channel.Get("0").Int()
				freq := channel.Get("1").String()
				ch <- prometheus.MustNewConstMetric(p.DSChannelPower, prometheus.GaugeValue, channel.Get("2").Float(), instance.Name, instance.Address, freq)
				ch <- prometheus.MustNewConstMetric(p.DSChannelSNR, prometheus.GaugeValue, channel.Get("3").Float(), instance.Name, instance.Address, freq)
				//modulation := channel.Get("4").String()
				status := channel.Get("5").String()
				if status == "Locked" {
					ch <- prometheus.MustNewConstMetric(p.DSChannelLocked, prometheus.GaugeValue, float64(1), instance.Name, instance.Address, freq)
				} else {
					ch <- prometheus.MustNewConstMetric(p.DSChannelLocked, prometheus.GaugeValue, float64(0), instance.Name, instance.Address, freq)
				}

				ch <- prometheus.MustNewConstMetric(p.DSChannelRXMer, prometheus.GaugeValue, channel.Get("6").Float(), instance.Name, instance.Address, freq)
				ch <- prometheus.MustNewConstMetric(p.DSChannelPreRS, prometheus.GaugeValue, channel.Get("7").Float(), instance.Name, instance.Address, freq)
				ch <- prometheus.MustNewConstMetric(p.DSChannelPostRS, prometheus.GaugeValue, channel.Get("8").Float(), instance.Name, instance.Address, freq)

				//prerserrors := channel.Get("7").Int()
				//postrserrors := channel.Get("8").Int()
				////fmt.Printf("Docsis 3.0 DS Channel: %d @ %dHz, %f dBmV, SNR %f dB, Modulation %s, Status %s, RXMer %f dB, PreRS %d, Post RS %d\n",
				//	id,freq,power,snr,modulation,status,rxmer,prerserrors,postrserrors)
				return true
			})

			// 21 - US Channel - JSON
			usChannels := gjson.Parse(gjson.Get(string(body), "21").String())
			usChannels.ForEach(func(key, value gjson.Result) bool {
				channel := gjson.Parse(value.String())
				id := channel.Get("0").Int()
				if id != 0 {
					// 	USChannelPower *prometheus.Desc
					//	USChannelSymbolRate* prometheus.Desc
					//	USChannelTimeouts *prometheus.Desc

					freq := channel.Get("1").String()
					ch <- prometheus.MustNewConstMetric(p.USChannelPower, prometheus.GaugeValue, channel.Get("2").Float(), instance.Name, instance.Address, freq)
					//power := channel.Get("2").Float()
					//symbolrate := channel.Get("3").String()
					//modulation := channel.Get("4").String()
					//channeltype := channel.Get("5").String()
					ch <- prometheus.MustNewConstMetric(p.USChannelTimeouts, prometheus.GaugeValue, channel.Get("6").Float(), instance.Name, instance.Address, freq, "1")
					ch <- prometheus.MustNewConstMetric(p.USChannelTimeouts, prometheus.GaugeValue, channel.Get("7").Float(), instance.Name, instance.Address, freq, "2")
					ch <- prometheus.MustNewConstMetric(p.USChannelTimeouts, prometheus.GaugeValue, channel.Get("8").Float(), instance.Name, instance.Address, freq, "3")
					ch <- prometheus.MustNewConstMetric(p.USChannelTimeouts, prometheus.GaugeValue, channel.Get("9").Float(), instance.Name, instance.Address, freq, "4")
					//t1timeouts := channel.Get("6").Float()
					//t2timeouts := channel.Get("7").Float()
					//t3timeouts := channel.Get("8").Float()
					//t4timeouts := channel.Get("9").Float()
					////fmt.Printf("Docsis 3.0 US Channel: %d @ %d Hz, Power %f dBmV, Symbol Rate %s, Modulation %s, Channel Type %s, Timeouts T1:%d T2:%d T3:%d T4:%d\n",
					//	id, frequency, power, symbolrate, modulation, channeltype, t1timeouts, t2timeouts, t3timeouts, t4timeouts)
				}
				return true
			})
			// 22 - Network Data - JSON
			// 23 - 3.1 DS Channel - JSON
			dsChannels31 := gjson.Parse(gjson.Get(string(body), "23").String())
			dsChannels31.ForEach(func(key, value gjson.Result) bool {
				channel := gjson.Parse(value.String())
				id := channel.Get("0").String()
				//channelwidth := channel.Get("1").Int()
				//fft := channel.Get("2").String()
				//subcarriers := channel.Get("3").Int()
				//modulation := channel.Get("4").String()
				//firstsubcarrier := channel.Get("5").Int()
				//lockstatus := channel.Get("6").String()
				//rxmer := channel.Get("7").Float()
				//plcpower := channel.Get("8").Float()
				////prerserrors := channel.Get("9").Int()
				////postrserrors := channel.Get("10").Int()

				//freq := channel.Get("1").String()


				status := channel.Get("6").String()
				if status == "Locked" {
					ch <- prometheus.MustNewConstMetric(p.DS31ChannelLocked, prometheus.GaugeValue, float64(1), instance.Name, instance.Address, id)
				} else {
					ch <- prometheus.MustNewConstMetric(p.DS31ChannelLocked, prometheus.GaugeValue, float64(0), instance.Name, instance.Address, id)
				}
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelPLCPower, prometheus.GaugeValue, channel.Get("8").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelRXMer, prometheus.GaugeValue, channel.Get("6").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelPreRS, prometheus.GaugeValue, channel.Get("9").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelPostRS, prometheus.GaugeValue, channel.Get("10").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelFirstSubcarrier, prometheus.GaugeValue, channel.Get("5").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelSubcarriers, prometheus.GaugeValue, channel.Get("3").Float(), instance.Name, instance.Address, id)
				ch <- prometheus.MustNewConstMetric(p.DS31ChannelWidth, prometheus.GaugeValue, channel.Get("1").Float(), instance.Name, instance.Address, id)
				//fmt.Printf("Docsis 3.1 DS Channel: %d @ %d MHz Wide, FFT %s, Subcarriers %d, Modulation: %s, " +
				//	"First Subcarrier %d Hz, Lock Status %s, RXMer %f dB, PLC Power %f dBmV, Pre RS %d, Post RS %d\n",
				//	id, channelwidth, fft, subcarriers, modulation, firstsubcarrier, lockstatus, rxmer, plcpower, prerserrors, postrserrors)
				return true
			})
			// 24 - 3.1 US Channel - JSON
			//usChannels31 := gjson.Parse(gjson.Get(string(body), "24").String())
			//usChannels31.ForEach(func(key, value gjson.Result) bool {
			//	channel := gjson.Parse(value.String())
			//	id := channel.Get("0").Int()
			//	if id != 0 {
			//		frequency := channel.Get("1").Int()
			//
			//		power := channel.Get("2").Float()
			//		symbolrate := channel.Get("3").String()
			//		modulation := channel.Get("4").String()
			//		channeltype := channel.Get("5").String()
			//		t1timeouts := channel.Get("6").Int()
			//		t2timeouts := channel.Get("7").Int()
			//		t3timeouts := channel.Get("8").Int()
			//		t4timeouts := channel.Get("9").Int()
			//		fmt.Printf("Docsis 3.1 US Channel: %d @ %d Hz, Power %f dBmV, Symbol Rate %s, Modulation %s, Channel Type %s, Timeouts T1:%d T2:%d T3:%d T4:%d\n",
			//			id, frequency, power, symbolrate, modulation, channeltype, t1timeouts, t2timeouts, t3timeouts, t4timeouts)
			//	}
			//	return true
			//})
			// 25 - US Number
			value = gjson.Get(string(body), "25")
			ch <- prometheus.MustNewConstMetric(p.USNumber, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 26 - DS Number
			value = gjson.Get(string(body), "26")
			ch <- prometheus.MustNewConstMetric(p.DSNumber, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)

			// 27 - 3.1 US Number
			value = gjson.Get(string(body), "27")
			ch <- prometheus.MustNewConstMetric(p.USNumber31, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)
			// 28 - 3.1 DS Number
			value = gjson.Get(string(body), "28")
			ch <- prometheus.MustNewConstMetric(p.DSNumber31, prometheus.GaugeValue, value.Float(), instance.Name, instance.Address)


			// 29 - Primary Channel type
			//value = gjson.Get(string(body), "29")
			//fmt.Printf("Primary Channel Type: %s\n", value)

			instanceWG.Done()
		}(instance)
	}
	// Wait for all instances to complete their poll
	instanceWG.Wait()
}