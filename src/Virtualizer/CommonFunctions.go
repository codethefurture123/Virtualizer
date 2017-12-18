package main

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/clbanning/mxj"
	"math/rand"
	crand "crypto/rand" 
    "time"
    //"os"
    "errors"
      "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
   // "reflect"
)
      
func tagextractor(data []byte, tag string) string{
	m, err := mxj.NewMapXml(data)
	
	if err != nil {
		//fmt.Println("NewMapXml err:", err)
	}
	xval,xerr:=m.ValueForPath(tag)
	if xerr!=nil{
	//fmt.Println("value not found")
	}
//	fmt.Println(xval)
	var tagval string
	switch xval.(type){
	 	case string: tagval=xval.(string)
	 	case map[string]interface{}: tagval=xval.(map[string]interface{})["#text"].(string)
	 	default: tagval=""
	}
	if tagval==""{
		return tag
	}else{
		return tagval
	}
}

func addDelay(delay time.Duration,ch chan<- bool){
	time.Sleep(delay * time.Second)
//	fmt.Println("woke up")
	ch <- true
}

func getRandomNumber(min int, max int) int{
	rand.Seed(time.Now().UnixNano())
	rnum:=rand.Intn(max-min)+min
	//fmt.Println(rnum)
    return rnum
}

func getFormattedTimeStamp(format string) string{
	ts:=time.Now().Format(format)//20060102150405.00
	return ts
}

func getFormattedTimeStampWithOffset(format string, d time.Duration, unit string ) string{
	
		var ts string		
	switch unit{
		case "s": ts=time.Now().Add(d*time.Second).Format(format)
		case "m": ts=time.Now().Add(d*time.Minute).Format(format)
		case "h": ts=time.Now().Add(d*time.Hour).Format(format)
		default : panic(errors.New("unit not recognized pass s for seconds m for minutes or h for hour"))
	}
	
		return ts	
	
}

func shuffle(data []byte,source string) string{
	if strings.Contains(source,"."){
		dest:=tagextractor(data , source)
		 return string(dest)
	}else
	{
	 rand.Seed(time.Now().UnixNano())
     src:=[]byte(source)
     dest := make([]byte, len(src))
     perm := rand.Perm(len(src))
     for i, v := range perm {
	     dest[v] = src[i]
     }
     return string(dest)
	}
	
}

func getGUID() string{
	// generate 32 bits timestamp
 	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := crand.Read(buff)
	
 	if numRead != len(buff) || err != nil {
 		panic(err)
 	}

 	uuid:=fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
 	return uuid
}



func DBInsertValue(column string,value string,dbName string, collectionName string) string{
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	//c := session.DB("test").C("CafeDB")
	c := session.DB(dbName).C(collectionName)
	
	err = c.Insert(bson.M {column: value})
	if err != nil {
		panic(err)
	}
	return "Success"
}

		
//Function to fetch a value from DB decided by the a value in the request
	
	func DBFetch(fetchDetails string,requestValue string,matchColumnName string) string{
	
	z := strings.Split(fetchDetails, ",")
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(z[1]).C(z[2])

	//fmt.Println("Macthing value",requestValue)
	var m bson.M
	err = c.Find(bson.M{matchColumnName: requestValue}).One(&m)
	//fmt.Println("Search Names:",matchColumnName)
	//panic(err)
	
	for key, value := range m{
		if key ==z[0]{//comparing the DB key whose value will match the value sent in z[0]

    return value.(string)
		}
	}
	return "Value Not Found"
}   
	
// Function to do dynamic correlation,i.e. to construct the response dynamically decided by the number of occurence of a 
//tag in the incoming request
 	
	func tagextractorForArrayWithCorrelation(data []byte,tag string, repeatStringFromReq string, occurence int, tagtocapture string )string {
		m, err := mxj.NewMapXml(data)
	   if err != nil {
			
	}
	   
	
    values, _ := m.ValuesForPath(repeatStringFromReq+".*")
//values, _ := m.ValuesForPath(dynamicValueToReplace+".*")
    
    var str string
    
    switch values[occurence].(type) {
 	case string: str, _ = values[occurence].(string)
 	case map[string]interface{}: str=values[occurence].(map[string]interface{})[tagtocapture].(string)
 	default: str="could not found"
 }
    
    
   	
   /* repeatTagArray:=strings.Split(tag, "Repeat(")
    repeatTag:=strings.Split(repeatTagArray[1], ")")*/
   /* returnTag:=strings.Replace(repeatString,"@@"+dynamicValueToReplace+"@@", str, -1)
      fmt.Println("dynamic ",returnTag)*/
			    		
    //return returnTag
    return str
   
}
	
	//Function to map the incoming request to a response or multiple responses.The call is frmo HandelerCreater.go
	func MultipleResponses(operation Operation,reqbody string,hah []byte,response string) string{
	for _,output := range operation.Outputs {
	          	//logrus.WithFields(logrus.Fields{"Output":output.Response}).Debug()
	          	
	          	
	          	correlationValue:=output.Tagvalue
	          	//fmt.Println("Before correlationValue:",correlationValue)
	          	        	
	          	if(strings.Contains(reqbody, correlationValue)){
		          	 response=output.Response
			   }else 
	          	if(strings.Contains(correlationValue,".")&&!strings.Contains(correlationValue,"=")){
	          			m,_ := mxj.NewMapXml(hah)
	          			_,xerr:=m.ValueForPath(correlationValue)
	          			//fmt.Println("value found in .",xval)
	          			if xerr==nil{
						//fmt.Println("value found")
						response=output.Response
						} else {
							 response="Sorry!,This request is not properly mapped to a response. Please check if the server is configured with a proper Xpath."
						}

	          		}else 
	          	if(strings.Contains(correlationValue,"=")){
	          			correlationTagValue:=strings.Split(correlationValue, "=")
	          			m,_ := mxj.NewMapXml(hah)
	          			xval,xerr:=m.ValueForPath(correlationTagValue[0])
	          			//fmt.Println("value found in =",xval)
	          			if xerr!=nil{
		          			 //fmt.Println("value not found")
	          			}
						if(xval==correlationTagValue[1]){
						response=output.Response
						} else {
							response="Sorry!,This request is not properly mapped to a response. Please check if the server is configured with a proper Xpath and condition value."
						}
						}
	          	//evaluatedIPVariables=evaluateInputVariables(output.Variables,hah)
			}
	          	return response
}

func Extract(tagValue string,delimiter string,index string,tailer string) string{
	 
	
		/*	if (delimiter=="comma"){
			 	valueArray:=strings.Split(tagValue, ",")
			 }else{ */
				valueArray:=strings.Split(tagValue, delimiter)
			 //}
			
			 
						i,err:= strconv.Atoi(index)
						 if err != nil {
						 	fmt.Println(err)
					     }
						extractedValue:=valueArray[i]
						if(tailer!="nil" &&tailer!="comma"){
							extractedValue=strings.TrimRight(extractedValue, tailer)
							
							}
						if(tailer=="comma"){
							extractedValue=strings.TrimRight(extractedValue,",")
							
							}
	return string(extractedValue)
}

func reserveItem(data []byte) string{


m, err := mxj.NewMapXml(data)
						if err != nil {
							panic(err)	
							}
						value, _ := m.ValuesForPath("Envelope.Body.reserveItem.request.telephoneNumberBlockReservationRequests.telephoneNumberBlockReservationRequest.blockCount.*")
						values, _ := m.ValuesForPath("Envelope.Body.reserveItem.request.telephoneNumberBlockReservationRequests.telephoneNumberBlockReservationRequest")

						blockCount,_:=m.ValuesForPath("Envelope.Body.reserveItem.request.telephoneNumberBlockReservationRequests.telephoneNumberBlockReservationRequest.blockCount" )
						
		    			groupID,_:=m.ValuesForPath("Envelope.Body.reserveItem.request.telephoneNumberBlockReservationRequests.telephoneNumberBlockReservationRequest.groupId")

finalString:=""
repeatString:="<typ:item xsi:type=\"typ:TelephoneNumberBlockItem\"><typ:location xsi:nil=\"true\"/><typ:reservedForCustomer xsi:nil=\"true\"/><typ:type xsi:nil=\"true\"/><typ:assignedToCustomer xsi:nil=\"true\"/><typ:market xsi:nil=\"true\"/><typ:configuration xsi:nil=\"true\"/><typ:id xsi:nil=\"true\"/><typ:clliCode xsi:nil=\"true\"/><typ:ported xsi:nil=\"true\"/><typ:blockCount>P_BLOCKCOUNT</typ:blockCount><typ:startTelephoneNumber>P_STN</typ:startTelephoneNumber><typ:endTelephoneNumber>P_ETN</typ:endTelephoneNumber><typ:portEligible xsi:nil=\"true\"/><typ:exceptions xsi:nil=\"true\"/><typ:disconnectReason xsi:nil=\"true\"/><typ:groupId>P_GROUPID</typ:groupId><typ:telephoneNumbers xsi:nil=\"true\"/></typ:item>"
//value :=tagextractor(data , "Envelope.Body.reserveItem.request.telephoneNumberBlockReservationRequests.telephoneNumberBlockReservationRequest.*" )
fmt.Println("Value:",value)
fmt.Println("Values:",values)

			    			fmt.Println("blockCounts:",blockCount)
			    			fmt.Println("groupID:",groupID)
			    			fmt.Println("value:",value)
occurence:=len(value)
fmt.Println("occurence:",occurence)
for i:=1;i<occurence;i++{
			    			fmt.Println("aArrays:",blockCount[i-1],groupID[i-1])
			    			finalString=repeatString+finalString
			    			/*strings.Replace(finalString, "P_BLOCKCOUNT",string(blockCount[i-1]), -1)
			    			strings.Replace(finalString, "P_STN", "9991", -1)
			    			strings.Replace(finalString, "P_ETN", "8881", -1)
			    			strings.Replace(finalString, "P_GROUPID", string(groupID[i-1]), -1)*/
			    			fmt.Println("final string:",i,finalString)
			    			 
			    			  
			    			/* for j:=1;j<=len(dynamicCorelation)-2;j=j+2{
			    			 	//fmt.Println("on chnage",NewrepeatString)
						//dynamicValueToReplace:=strings.Split(dynamicCorelation[j], "@@")
						dynamicValueToReplace:=dynamicCorelation[j]
						lastTagNames:=strings.Split(dynamicCorelation[j], ".")
						length:=len(lastTagNames) 
						lastTagName:=lastTagNames[length-1] //To get the last tag name
						str:=tagextractorForArrayWithCorrelation(data ,ip,repeatStringFromReq ,i, lastTagName)
						//fmt.Println("Value:",str)
						repeatString=strings.Replace(repeatString,"@@"+dynamicValueToReplace+"@@", str, 1)
						*/
	  
						
						}
			    			 
			    			 return finalString
			    			  


}

