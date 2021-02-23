package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"

	//"github.com/StackExchange/wmi"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type SysUtil struct{}

/**
获取本地IP地址
*/
func GetLocalIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type ModelType struct {
	typeName    string
	typeNameAll string
	TypeOf      *reflect.Type
}

func GetInterface(c interface{}) {
	str := reflect.TypeOf(c).String()
	i := strings.LastIndex(str, ".") + 1
	fmt.Println(str[i:])
}
func GetModelType(model interface{}) ModelType {
	modelType := reflect.TypeOf(model)
	mt := ModelType{}
	mt.TypeOf = &modelType
	mt.typeNameAll = modelType.String()
	i := strings.LastIndex(mt.typeNameAll, ".")
	if i < 0 {
		mt.typeName = mt.typeNameAll
	} else {
		mt.typeName = mt.typeNameAll[i+1:]
	}
	return mt
}
func GetModelTypeName(modelType reflect.Type) string {
	mt := ModelType{}
	mt.TypeOf = &modelType
	mt.typeNameAll = modelType.String()
	i := strings.LastIndex(mt.typeNameAll, ".")
	if i < 0 {
		mt.typeName = mt.typeNameAll
	} else {
		mt.typeName = mt.typeNameAll[i+1:]
	}
	return mt.typeName
}

//字符串判断是否为空
func IsNoBlank(model interface{}) bool {
	if model == nil {
		return false
	}
	return !IsBlank(model)
}

func IsBlankValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		values := strings.Replace(value.String(), " ", "", -1)
		// 去除换行符
		values = strings.Replace(values, "\n", "", -1)
		values = strings.Replace(values, "\t", "", -1)
		return len(values) == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
func IsBlank(model interface{}) bool {
	if model == nil {
		return true
	}
	return IsBlankValue(reflect.ValueOf(model))
}

type SystemMac struct {
	NetInterfaces []net.Interface
	Addrs         []net.Addr
	Ip            string
}

func (s *SystemMac) GetSystemMac() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err.Error())
	}
	s.NetInterfaces = netInterfaces
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			s.Addrs = addrs
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						s.Ip = ipnet.IP.String()
					}
				}
			}
		}
	}
}

type intfInfo struct {
	Name string
	Ipv4 []string
	Ipv6 []string
}

//网卡信息
func GetIntfs() []intfInfo {
	intf, err := net.Interfaces()
	if err != nil {
		return []intfInfo{}
	}
	var is = make([]intfInfo, len(intf))
	for i, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		is[i].Name = v.Name
		for _, ip := range ips {
			if strings.Contains(ip.String(), ":") {
				is[i].Ipv6 = append(is[i].Ipv6, ip.String())
			} else {
				is[i].Ipv4 = append(is[i].Ipv4, ip.String())
			}
		}
	}
	return is
}

type TryCatch struct {
	errChan      chan interface{}
	catches      map[reflect.Type]func(err error)
	defaultCatch func(err error)
}

func (t TryCatch) Try(block func()) TryCatch {
	t.errChan = make(chan interface{})
	t.catches = map[reflect.Type]func(err error){}
	t.defaultCatch = func(err error) {}
	go func() {
		defer func() {
			t.errChan <- recover()
		}()
		block()
	}()
	return t
}

func (t TryCatch) CatchAll(block func(err error)) TryCatch {
	t.defaultCatch = block
	return t
}

func (t TryCatch) Catch(e error, block func(err error)) TryCatch {
	errorType := reflect.TypeOf(e)
	t.catches[errorType] = block
	return t
}

func (t TryCatch) Finally(block func()) TryCatch {
	err := <-t.errChan
	if err != nil {
		catch := t.catches[reflect.TypeOf(err)]
		if catch != nil {
			catch(err.(error))
		} else {
			t.defaultCatch(err.(error))
		}
	}
	block()
	return t
}

type OSoftware struct {
	Id      string `json:"id"` // id String 主键
	Version string `json:"version"`
	Md5     string `json:"md5"`
	Data    []byte `json:"data"`
	Size    int    `json:"size"`
	Code    string `json:"code"`
}

func Cmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	buf, _ := cmd.Output()
	fmt.Println(string(buf))
}

// 根据IP地址查询MAC地址
func GetByIpToMac(ip string) (mac string, err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("arp", "-a", ip)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("error", err)
		return
	}
	mac = stdout.String()
	if len(mac) > 140 {
		mac = mac[strings.Index(mac, ip):]
		mac = mac[len(ip):]
		mac = strings.ReplaceAll(mac, " ", "")
		mac = strings.ToUpper(mac[:17])
	}
	return
}

func ExeFun(name string, arg ...string) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println(name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		if err.Error() != "exit status 1" {
			fmt.Printf("error: %v\n", err)
		}
	}
	str := stderr.String()
	if len(str) == 0 {
		str = stdout.String()
	}
	return str
}

//region 执行命令
func ExeFun2(name string, arg ...string) string {
	cmd0 := exec.Command(name, arg...)
	stdout, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Println("获取通道失败：", err)
		return "500"
	}
	if err := cmd0.Start(); err != nil {
		fmt.Println("命令执行失败：", err)
		return "501"
	}
	var outputBufo bytes.Buffer
	for {
		output0 := make([]byte, 5)
		n, err := stdout.Read(output0)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error:找不到通道", err)
				return "502"

			}
		}
		if n > 0 {
			outputBufo.Write(output0[:n])
		}
	}
	str := outputBufo.String()
	return str
}

func ExeFunErr(name string, arg ...string) (str string, err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println(name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println(err.Error())
		return
	}
	if stderr.String() != "" {
		str = stderr.String()
	} else {
		str = stdout.String()
	}
	return
}

//endregion

//ip到数字
func Ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

//数字到IP
func BackToIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

//获取本机Mac
func GetMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Get loacl Mac failed")
		return "", nil
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr
		if mac.String() != "" {
			return mac.String(), nil
		}
	}
	return "", nil
}
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
