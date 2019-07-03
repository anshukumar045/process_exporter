// +build !noprocess

package collector

import (
        "log"
        "io/ioutil"
                "strconv"
                "fmt"

        "github.com/prometheus/client_golang/prometheus"
//      "github.com/prometheus/common/log"
)

type processCollector1 struct {
    smapsinfo          *prometheus.Desc
        totalvm            *prometheus.Desc
    statusinfo         *prometheus.Desc
    cpu_info           *prometheus.Desc
    fileinfo           *prometheus.Desc
        incounterinfo      *prometheus.Desc
}

func init() {
        registerCollector("process", defaultEnabled, NewProcessCollector)
}

const (
        processCollectorSubsystem = "process"
)

// NewProcessCollector returns a new Collector exposing kernel/system statistics.
func NewProcessCollector() (Collector, error) {
        return &processCollector1{
                smapsinfo: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "mem_rss"),
                        "rss",
                        []string{"pid","cmdline","type"}, nil,
                ),
                 totalvm: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "virtual_memory"),
                        "total memory",
                        []string{"type"},nil,
                ),
                statusinfo: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "status_info"),
                        "Status info",
                        []string{"pid","cmd_line","type"}, nil,
                ),
                fileinfo: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "files_info"),
                        "Files Info",
                        []string{"type"}, nil,
                ),
                cpu_info: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "cpu_info"),
                        "Cpuinfo user system nice iowait guest Minpgflt Majpgflt ThreadCnt",
                        []string{"pid","cmd_line","mode"}, nil,
                ),
                incounterinfo: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, processCollectorSubsystem, "counter_info"),
                        "Read Count",
                        []string{"pid","cmd_line","type"}, nil,
                ),
        },nil

}

// Update implements Collector and exposes process related metrics from /proc/stat and /sys/.../process/.
func (c *processCollector1) Update(ch chan<- prometheus.Metric) error {
                // status := PrintCache("collector")
                // if status == false {
                //      fmt.Println("No value found")
                // }
        if err := c.processStat(ch); err != nil {
                return err
        }
        return nil
}


// updateStat reads /proc/stat through procfs and exports process related metrics.
func (c *processCollector1) processStat(ch chan<- prometheus.Metric) error {

                if x, found := C.Get("metricconf"); found {
                        foo := x.(*[]*Config)
                                        for _, vl := range *foo {
                                                if vl.Typ == "process" {
                                                        fmt.Println(vl)
                                                }
                        }
                }

                filelst, err := ioutil.ReadDir("/proc")
                        if err !=nil {
                                log.Println(err)
                                }

                for _, f := range filelst {
                pid, err := strconv.Atoi(f.Name())
                        if err == nil {
                                cmdline , err := getCommandline(pid)   // command line for each PID
                                      if err != nil {
                        log.Println(err)
                      }
                                if len(cmdline) > 0 {
                                mycppu, err :=CPUinfo(pid)            // cpu-total, user, sys, nice, iowait, guest
                                      if err != nil {
                        log.Println(err)
                      }
//                              myiocnt, err := getIOCounter(pid)       // readcount, write count, read_byte, write_byte
//                                    if err != nil {
//                       log.Println(err)
//                     }
//                              mysmaps, err := getSmaps(pid)               // rss, pss, shared dirty, private dirty
//                                    if err != nil {
//                       log.Println(err)
//                    }
                                mystatusinfo, err := getStatusInfo(pid)               // vmdata, vmlck, vmstk
                                      if err != nil {
                        log.Println(err)
                      }

/********************************************************************************************************************************/
                                                        // STATS

                                        ch <- prometheus.MustNewConstMetric(
                                          c.cpu_info,
                                          prometheus.CounterValue,
                                          float64(mycppu.CPU_total),
                                          strconv.Itoa(pid),
                                          cmdline,
                                          "TotalUsage",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                          c.cpu_info,
                                          prometheus.CounterValue,
                                          float64(mycppu.User),
                                          strconv.Itoa(pid),
                                          cmdline,
                                          "User",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,
                                           prometheus.CounterValue,
                                           float64(mycppu.System),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "System",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,
                                           prometheus.CounterValue,
                                           float64(mycppu.Nice),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "Nice",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,
                                           prometheus.CounterValue,
                                           float64(mycppu.Iowait),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "Iowait",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,
                                           prometheus.CounterValue,
                                           float64(mycppu.Guest),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "Guest",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,                                     // pgflt
                                           prometheus.CounterValue,
                                           float64(float64(mycppu.Minpgflt)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "MinorPgfalut",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,                                           // pgflt
                                           prometheus.CounterValue,
                                           float64(float64(mycppu.Majpgflt)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "MajorPgfalut",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.cpu_info,                                           //thrdcnt
                                           prometheus.CounterValue,
                                           float64(mycppu.ThreadCnt),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "ThreadCount",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.smapsinfo,
                                           prometheus.CounterValue,
                                           float64(1.00),
//                                         float64(float64(mysmaps.rss)/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "RSS",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.smapsinfo,
                                           prometheus.CounterValue,
                                           float64(1.00),
//                                         float64(float64(mysmaps.pss)/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "PSS",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.smapsinfo,
                                           prometheus.CounterValue,
                                           float64(1.00),
//                                         float64(mysmaps.pdirty),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "PrivateDirty",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.statusinfo,
                                           prometheus.CounterValue,
                                           float64(float64(mystatusinfo.VmLck)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "VirtualMemoryLock",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.statusinfo,
                                           prometheus.CounterValue,
                                           float64(float64(mystatusinfo.VmStk)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "VirtualMemoryStack",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.statusinfo,
                                           prometheus.CounterValue,
                                           float64(float64(mystatusinfo.VmData)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "VirtualMemoryData",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.statusinfo,
                                           prometheus.CounterValue,
                                           float64(float64(mystatusinfo.VmExe)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "VirtualMemoryExec",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.statusinfo,
                                           prometheus.CounterValue,
                                           float64(float64(mystatusinfo.VmLib)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "VirtualMemoryLib",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.incounterinfo,
                                           prometheus.CounterValue,
                                           float64(1.00),
//                                         float64(myiocnt.rdcnt),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "ReadCount",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                            c.incounterinfo,
                                            prometheus.CounterValue,
                                            float64(1.00),
//                                          float64(myiocnt.wrtcnt),
                                            strconv.Itoa(pid),
                                            cmdline,
                                            "WriteCount",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.incounterinfo,
                                           prometheus.CounterValue,
                                           float64(1.00),
//                                         float64(float64(myiocnt.rdbytes)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "ReadByte",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.incounterinfo,
                                           prometheus.CounterValue,
                                           float64(1),
//                                         float64(float64(myiocnt.wrtbytes)/1024/1024),
                                           strconv.Itoa(pid),
                                           cmdline,
                                           "WriteByte",
                                        )
                                        }
                                }
                        }

                                        myfilestat, err := getFileStat()
                                            if err != nil {
                                                log.Println(err)
                                             }

                                        ch <- prometheus.MustNewConstMetric(
                                           c.fileinfo,
                                           prometheus.CounterValue,
                                           float64(myfilestat.OpenFiles),
                                           "OpenFiles",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.fileinfo,
                                           prometheus.CounterValue,
                                           float64(myfilestat.TotalFiles),
                                           "TotalFiles",
                                        )

                                        mymeminfo, err := getMemInfo()
                                            if err!= nil {
                                                log.Println(err)
                                            }

                                        ch <- prometheus.MustNewConstMetric(
                                           c.totalvm,
                                           prometheus.CounterValue,
                                           float64(float64(mymeminfo.MemTotal)/1024/1024),
                                           "Vmem Total",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.totalvm,
                                           prometheus.CounterValue,
                                           float64(float64(mymeminfo.MemFree)/1024/1024),
                                           "VmemFree",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.totalvm,
                                           prometheus.CounterValue,
                                           float64(float64(mymeminfo.MemTotal)/1024/1024  -  float64(mymeminfo.MemFree)/1024/1024),
                                           "VmemUsed",
                                        )

                                        ch <- prometheus.MustNewConstMetric(
                                           c.totalvm,
                                           prometheus.CounterValue,
                                           float64(float64(mymeminfo.MemAvailable)/1024/1024),
                                           "VmemAvailable",
                                        )

 return nil

}
