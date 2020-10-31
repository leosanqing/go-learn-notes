package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/tidwall/sjson"
	"github.com/zhixinlian/zxl-go-sdk/sm/sm2"
	"math/big"
	"reflect"
	"strconv"
	"sync"
	"time"
)



const (
	CREDIT = "CREDIT"
	CASH   = "CASH"

	CHECKED = "CHECKED"

	DataReq = "channel1111_share_request_audit:"

	DataShare = "DATA_SHARE:"

	DataAudit = "DATA_AUDIT:"
)

type SimpleChainCode struct {
}

type DataShareBsi struct {

	// id
	Id int `json:"id,omitempty"`

	Uid int `json:"uid,omitempty"`

	ShareDataId int `json:"shareDataId,omitempty"`

	//TODO.
	ShareObjects []ShareLogicDetails `json:"shareLogicDetails,omitempty"`

	Desc string `json:"desc"`

	CreateTime time.Time `json:"createTime"`

	ChainHash string `json:"chainHash"`

	Sign string `json:"sign,omitempty"`

	Pk string `json:"pk,omitempty"`
}

type DataReadBsi struct {

}
type DataSharePO struct {

	// id
	Id int `json:"id,omitempty"`

	ShareDataId int `json:"shareDataId,omitempty"`

	ShareObjects []ShareObject `json:"shareObjects,omitempty"`

	Desc string `json:"desc"`

	AsymmetricKey string `json:"asymmetricKey"`

	CreateTime time.Time `json:"createTime"`

	ChainHash string `json:"chainHash"`
}

type ShareLogicDetails struct {
	Id int `json:"id"`

	Uid int `json:"uid,omitempty"`

	ShareType int `json:"shareType,omitempty"`

	ShareExpire int64 `json:"shareExpire,omitempty"`

	AsymmetricKey string `json:"asymmetricKey"`

	ShareLogic shareLogic `json:"shareLogic"`
}

type shareLogic struct {
	Id int `json:"id,omitempty"`

	ShareDataId int `json:"shareDataId,omitempty"`

	ShareLogicDesc string `json:"shareLogicDesc"`

	CreateTime int64 `json:"createTime"`

	ChainHash string `json:"chainHash"`
}

type ShareObject struct {
	Id int `json:"id"`

	Uid int `json:"uid,omitempty"`

	ShareType int `json:"shareType,omitempty"`

	ShareExpire time.Time `json:"shareExpire,omitempty"`

	AsymmetricKey string `json:"asymmetricKey"`

	ShareLogicId int `json:"shareLogicId"`

	CreateTime time.Time `json:"createTime"`
}

type DataReqBsi struct {
	Id int `json:"id"`

	ShareDataId int `json:"shareDataId,omitempty"`

	RequestUid int `json:"requestUid,omitempty"`

	Desc string `json:"desc"`

	AuditStatus int `json:"auditStatus,omitempty"`

	Sign string `json:"sign,omitempty"`

	CreateTime time.Time `json:"createTime"`

	ChainHash string `json:"chainHash"`

	AuditUid int `json:"auditUid,omitempty"`

	ShareObject string `json:"shareObject,omitempty"`

	Pk string `json:"pk,omitempty"`
}

type DataReqPO struct {
	Id int `json:"id"`

	ShareDataId int `json:"shareDataId,omitempty"`

	RequestUid int `json:"requestUid,omitempty"`

	Desc string `json:"desc"`

	AuditStatus int `json:"auditStatus,omitempty"`

	CreateTime time.Time `json:"createTime"`

	ChainHash string `json:"chainHash"`

	AuditUid int `json:"auditUid,omitempty"`

	ShareObjects string `json:"shareObject,omitempty"`
}

func main() {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	} else {
		fmt.Println("SimpleChainCode successfully started")
	}
}

func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Init")
	return shim.Success(nil)
}

func (t *SimpleChainCode) Query(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Query")
	return shim.Success(nil)
}

func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "dataShare":
		response := t.dataShare(stub, args)
		if response.Status == 200 {
			return shim.Success(nil)
		}
		return shim.Error(response.Message)
	case "dataReq":
		response := t.dataReq(stub, args)
		if response.Status == 200 {
			return shim.Success(nil)
		}
		return shim.Error(response.Message)
	case "dataAudit":
		response := t.dataAudit(stub, args)
		if response.Status == 200 {
			return shim.Success(nil)
		}
		return shim.Error(response.Message)
	}

	return shim.Error("Invalid invoke function name: " + function)
}

func (t *SimpleChainCode) dataShare(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	jsonString := args[0]
	var dataShare DataShareBsi
	err := json.Unmarshal([]byte(jsonString), &dataShare)
	if err != nil {
		return shim.Error("")
	}

	fmt.Println("数据共享上链请求数据：" + jsonString)

	// 验签
	buffer := bytes.Buffer{}

	// 生成需要验签的数据
	for _, object := range dataShare.ShareObjects {
		buffer.WriteString(strconv.Itoa(object.Uid))
		buffer.WriteString(strconv.Itoa(object.ShareType))
		buffer.WriteString(strconv.FormatInt(object.ShareExpire, 10))
	}
	// TODO 验签逻辑
	sign := fmt.Sprintf("%d%d%s%s", dataShare.Uid, dataShare.ShareDataId, buffer.String(), dataShare.Desc)
	r, _ := verify(dataShare.Pk, dataShare.Sign, sign)
	if !r {
		return shim.Error("Verify false of dataShare.")
	}

	newJson, _ := sjson.Delete(jsonString, "pk")
	fmt.Print(newJson)

	var dataSharePo DataSharePO
	SimpleCopyProperties(&dataShare, &dataSharePo)
	jsonStr, err := json.Marshal(dataSharePo)

	if err != nil {
		return shim.Error("Json 化对象失败" + strconv.Itoa(dataSharePo.Id))
	}

	err = stub.PutState(fmt.Sprintf("%s%d", DataShare, dataSharePo.Id), jsonStr)
	if err != nil {
		return shim.Error(" dataReqPo 上链失败" + strconv.Itoa(dataSharePo.Id))
	}
	return shim.Success(nil)
}

func (t *SimpleChainCode) dataReq(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	jsonString := args[0]
	var dataReqBsi DataReqBsi
	err := json.Unmarshal([]byte(jsonString), &dataReqBsi)
	if err != nil {
		return shim.Error("解析 dataReqBsi 失败")
	}

	fmt.Println("数据请求 上链请求数据：" + jsonString)

	// 验签
	sign := fmt.Sprintf("%d%d%s", dataReqBsi.RequestUid, dataReqBsi.ShareDataId, dataReqBsi.Desc)
	r, _ := verify(dataReqBsi.Pk, dataReqBsi.Sign, sign)
	if !r {
		return shim.Error("Verify false of dataReq.")
	}
	// BSI 对象转换为数据库对象
	//var dataReqPo DataReqPO
	//err = SimpleCopyProperties(&dataReqPo, &dataReqBsi)
	//if err != nil {
	//	return shim.Error("属性拷贝失败" + strconv.Itoa(dataReqPo.Id))
	//}
	//
	//var shareObj ShareObject
	//err = json.Unmarshal([]byte(dataReqBsi.ShareObject), &shareObj)
	//
	//if err != nil {
	//	return shim.Error("解析 shareObj 对象失败" + strconv.Itoa(dataReqPo.Id))
	//}
	//dataReqPo.ShareObjects = shareObj
	//jsonStr, err := json.Marshal(dataReqPo)
	//if err != nil {
	//	return shim.Error("Json 化对象失败" + strconv.Itoa(dataReqPo.Id))
	//}

	newJson, _ := removePkSignFromJsonStr(jsonString)
	if len(newJson) == 0 {
		return shim.Error("去除Pk Sign失败")
	}
	err = stub.PutState(fmt.Sprintf("%s%d", DataReq, dataReqBsi.Id), newJson)
	if err != nil {
		return shim.Error(" dataReqPo 上链失败" + strconv.Itoa(dataReqBsi.Id))
	}
	return shim.Success(nil)
}

func removePkSignFromJsonStr(oldJson string) ([]byte, error) {
	newJson, err := sjson.Delete(oldJson, "pk")
	if err != nil {
		return []byte(""), err
	}
	newJson, err = sjson.Delete(newJson, "sign")
	if err != nil {
		return []byte(""), err
	}
	return []byte(newJson), err
}

func (t *SimpleChainCode) dataAudit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	jsonString := args[0]
	var dataReqBsi DataReqBsi
	err := json.Unmarshal([]byte(jsonString), &dataReqBsi)
	if err != nil {
		return shim.Error("")
	}

	fmt.Println("数据审核 上链请求数据：" + jsonString)

	// 验签
	sign := fmt.Sprintf("%d%d%s", dataReqBsi.Id, dataReqBsi.AuditUid, dataReqBsi.Desc)
	r, _ := verify(dataReqBsi.Pk, dataReqBsi.Sign, sign)
	if !r {
		return shim.Error("Verify false of dataReq.")
	}

	newJson, _ := removePkSignFromJsonStr(jsonString)
	if len(newJson) == 0 {
		return shim.Error("去除Pk Sign失败")
	}
	//var dataReqPo DataReqPO
	//err = SimpleCopyProperties(&dataReqPo, &dataReqBsi)
	//jsonStr, err := json.Marshal(dataReqPo)
	//if err != nil {
	//	return shim.Error("Json 化对象失败" + strconv.Itoa(dataReqPo.Id))
	//}

	err = stub.PutState(fmt.Sprintf("%s%d", DataAudit, dataReqBsi.AuditUid), newJson)
	if err != nil {
		return shim.Error("dataReqPo 上链失败" + strconv.Itoa(dataReqBsi.Id))
	}

	return shim.Success(nil)
}
//
//func (t *SimpleChainCode) dataValidate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//	if len(args) != 1 {
//		return shim.Error("Incorrect number of arguments. Expecting 1")
//	}
//	jsonString := args[0]
//
//	var dataReadAction DataReadAction
//	err := json.Unmarshal([]byte(jsonString), &dataReadAction)
//
//	if err != nil {
//		return shim.Error("")
//	}
//
//	objectBytes, _ := stub.GetState("DataPub" + dataReadAction.dataId)
//
//	// 查看数据权限
//	var dataPub DataPub
//	err = json.Unmarshal([]byte(objectBytes), &dataPub)
//	if err != nil {
//		return shim.Error("")
//	}
//
//	//  校验分享者集合，有效期
//	objectBytes1, _ := stub.GetState(DATA_SHARE + dataReadAction.dataId)
//	var dataShare DataReqBsi
//	err = json.Unmarshal([]byte(objectBytes1), &dataShare)
//
//	array := dataShare.ShareObjects
//
//	for _, value := range array {
//		if value.Uid != dataReadAction.Uid {
//			continue
//		}
//		if value.ShareExpire.Before(time.Now()) {
//			return shim.Error("已超过有效时间，请重新申请相应权限")
//		}
//		if dataPub.data == nil && dataPub.address == nil {
//			return shim.Error("链上数据无内容")
//		}
//		return shim.Success(dataShare.AsymmetricKey)
//	}
//
//	return shim.Error("该用户没有权限查看数据")
//}


//定义签名算法
func verify(pkStr, sigStr, msgString string) (ok bool, err error) {
	pkBytes, err := hex.DecodeString(pkStr)
	if err != nil {
		return false, fmt.Errorf("incorrect pk format, %s", err)
	}
	sigBytes, err := hex.DecodeString(sigStr)
	if err != nil {
		return false, fmt.Errorf("incorrect signature format, %s", err)
	}
	msgBytes := []byte(msgString)
	x, y := elliptic.Unmarshal(curve, pkBytes)
	pubilcKey := &sm2.PublicKey{curve, x, y}
	sig := new(Signature)
	_, err = asn1.Unmarshal(sigBytes, sig)
	if err != nil {
		return false, err
	}

	result := sm2.Verify(pubilcKey, msgBytes, sig.R, sig.S)

	return result, nil
}

var curve = P256Sm2()

type p256Curve struct {
	*elliptic.CurveParams
}

var p256Sm2Params *elliptic.CurveParams
var p256sm2Curve p256Curve
var initonce sync.Once

type Signature struct {
	R, S *big.Int
}

// 取自elliptic的p256.go文件，修改曲线参数为sm2
// See FIPS 186-3, section D.2.3
func initP256Sm2() {
	p256Sm2Params = &elliptic.CurveParams{Name: "SM2-P-256"} // 注明为SM2
	//SM2椭 椭 圆 曲 线 公 钥 密 码 算 法 推 荐 曲 线 参 数
	p256Sm2Params.P, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	p256Sm2Params.N, _ = new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	p256Sm2Params.B, _ = new(big.Int).SetString("28E9FA9E9D9F5E344D5A9E4BCF6509A7F39789F515AB8F92DDBCBD414D940E93", 16)
	p256Sm2Params.Gx, _ = new(big.Int).SetString("32C4AE2C1F1981195F9904466A39C9948FE30BBFF2660BE1715A4589334C74C7", 16)
	p256Sm2Params.Gy, _ = new(big.Int).SetString("BC3736A2F4F6779C59BDCEE36B692153D0A9877CC62A474002DF32E52139F0A0", 16)
	p256Sm2Params.BitSize = 256

	p256sm2Curve = p256Curve{p256Sm2Params}
}

func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

func P256Sm2() elliptic.Curve {
	initonce.Do(initP256Sm2)
	return p256sm2Curve
}

