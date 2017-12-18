package main

import (
	
	"strconv"
	"strings"
	"time"
	"fmt"
	"github.com/clbanning/mxj"
	"bytes"
	"os/exec"
	 "io/ioutil"
	// "fmt"
	

)
func evaluateInputVariables(rawmap map[string]string, data []byte) map[string]string{
 //fmt.Println("inside IP evaluation")

 //var evaluatedIPVariablesForRepeat map[string]string

 varmap := make(map[string]string)
	for k,ip:= range rawmap{
		
	
					//Common function :1
			    	varmap[k]=ip
			    	if strings.Contains(ip, "getRandomNumber") {
				    	z := strings.Split(ip, "(")
				    	z[1]=strings.TrimRight(z[1], ")")
				    	limit:=strings.Split(z[1],",")
				    	min,_:=strconv.Atoi(limit[0])
				    	max,_:=strconv.Atoi(limit[1])
					    varmap[k]=strconv.Itoa(getRandomNumber(min,max))
			    		
			    	}else//Common function :2
			    	if strings.Contains(ip,"getFormattedTimeStampWithOffset"){
			    		
			    		z := strings.Split(ip, "(")
			    		z[1]=strings.TrimRight(z[1], ")")
			    		y:=strings.Split(z[1],",")
			    	
			    		duration,_:=strconv.Atoi(y[1])
			     		varmap[k]=getFormattedTimeStampWithOffset(y[0], time.Duration(duration), y[2])

			    	}else//Common function :3
			    	if strings.Contains(ip, "getFormattedTimeStamp"){
			    		z := strings.Split(ip, "(")
				    	z[1]=strings.TrimRight(z[1], ")")
				    	varmap[k]=getFormattedTimeStamp(z[1])
				    	
					}else//Common function :4
			    	if strings.Contains(ip, "shuffle"){
			    		z := strings.Split(ip, "(")
				    	z[1]=strings.TrimRight(z[1], ")")
			    		varmap[k]=shuffle(data,z[1])
			    	}else//Common function :5
			    	if strings.Contains(ip,"getGUID"){
			    		varmap[k]=getGUID()
			    	}else//Common function :6
			    	if strings.Contains(ip,"DBInsertValue"){
			    		z := strings.Split(ip, "(")
				    	z[1]=strings.TrimRight(z[1], ")")
				    	y:=strings.Split(z[1], ",")
			    		value :=tagextractor(data, y[0])
			    		varmap[k]=DBInsertValue(k,value,y[1],y[2])
			    	}else//Common function :7
			    	if strings.Contains(ip,"DBFetch"){
			    		z := strings.Split(ip, "(")
				    	y := strings.Split(z[1], ")") //y[0]
				    	x := strings.Split(ip, "when(")
				    	w := strings.Split(x[1], ")") //w[0]
				    	value :=tagextractor(data, w[0])
				    	v := strings.Split(ip, "matches(")
				    	u := strings.Split(v[1], ")") //u[0]
			    		
			    		varmap[k]=DBFetch(y[0],value,u[0])
			    	}else//Common function :9 For fetching specific value from request
			    	if strings.Contains(ip,"Extract"){
			    		
			    		z := strings.Split(ip, "(")
				    	y := strings.Split(z[1], ")") //y[0]
						w := strings.Split(y[0], ",") //y[0]
						value :=tagextractor(data, w[0])
						
						varmap[k]=Extract(value, w[1],w[2],w[3])//(value,delimiter,index,tailer)
						//fmt.Println("Extract Caled",varmap[k])
					
						
			    	}else//Common function :10 For Calling external Java calls
			    	if strings.Contains(ip,"Java"){
			    		
			    		z := strings.Split(ip, "(")
				    	y := strings.Split(z[1], ")") //y[0]
				    	filePath:= strings.Split(y[0], ",")
				    	
				    	inputConvrt2String:=string(data[:])
				    	commandLineInput:=" "+"\""+ inputConvrt2String +"\""
				    	cmdToInvokeJave := "java -cp . "+filePath[0]+ commandLineInput
				    	//fmt.Println(cmdToInvokeJave)
						//command, err := exec.Command("bash", "-c", cmdToInvokeJave).Output()
						 exec.Command("bash", "-c", cmdToInvokeJave).Output()
				    	//fmt.Println("command executed, file is being read",command)
	                	 responseFromJava, err := ioutil.ReadFile(filePath[1])
						
				    	varmap[k]=string(responseFromJava)
				    	if err!=nil{
				    		panic(err)
				    	}
					
						
			    	}else //Common function :8
			    	if ( strings.Contains(ip,"Repeat") ){
				    	NewrepeatString:=""

			    		var returnString bytes.Buffer
			    		z := strings.Split(ip, "wrtTag(")
			    		repeatStringFromReq:=strings.TrimRight(z[1], ")")//WRT string
			    		
			    		 repeatTagArray:=strings.Split(ip, "Repeat(")
						 repeatTag:=strings.Split(repeatTagArray[1], ")")
			    		 repeatString:=repeatTag[0]
			    								
						m, err := mxj.NewMapXml(data)
						if err != nil {
							panic(err)	
							}
						values, _ := m.ValuesForPath(repeatStringFromReq+".*")//  Envelope.Body.AddBVEUids.addBVEUidInfo.bveUIDInfo.bveUid.*
						fmt.Println("values:",values)
						/*if {strings.Contains(values,)
							
						}*/
						//fmt.Println("Occurence, len:", len(values))
						occurence:=len(values)
						//fmt.Println("Occurence",occurence)
						
			    		if strings.Contains(ip,"@@"){
			    			 
			    			 //	fmt.Println("inside multiple @@ for",j)
			    		//dynamicCorelation:=strings.Split(ip, "@@")
			    		//dynamicCount:=strings.Count(ip, "@@")
			    		for i:=1;i<occurence;i++{
			    			 dynamicCorelation:=strings.Split(ip, "@@")
			    			  repeatString:=repeatTag[0]
			    			 for j:=1;j<=len(dynamicCorelation)-2;j=j+2{
			    			 	//fmt.Println("on chnage",NewrepeatString)
						//dynamicValueToReplace:=strings.Split(dynamicCorelation[j], "@@")
						dynamicValueToReplace:=dynamicCorelation[j]
						lastTagNames:=strings.Split(dynamicCorelation[j], ".")
						length:=len(lastTagNames) 
						lastTagName:=lastTagNames[length-1] //To get the last tag name
						str:=tagextractorForArrayWithCorrelation(data ,ip,repeatStringFromReq ,i, lastTagName)
						//fmt.Println("Value:",str)
						repeatString=strings.Replace(repeatString,"@@"+dynamicValueToReplace+"@@", str, 1)
						
	  
						
						}
			    			  NewrepeatString=NewrepeatString+repeatString
			    			  
			    			 
			    			 }
			    		startdelimiter:= "${"
						enddelimeter:="}"
			    		for strings.Contains(NewrepeatString, startdelimiter){
             
							       z := strings.SplitN(NewrepeatString, startdelimiter,2)
						         y := strings.SplitN(z[1], enddelimeter,2)
						        // fmt.Println("response splits...",evaluatedIPVariables[y[0]])
						        // fmt.Println("response splits...",y[0])
						         NewrepeatString=z[0]+varmap[y[0]]+y[1]
						         fmt.Println("Updated Repeat String:",NewrepeatString)
   
								}
			varmap[k]=NewrepeatString
			    	
			    	}else {
			    		for i:=1;i<occurence;i++{	  
						returnString.WriteString(repeatString)
						}
			    	}
			    		
//			    		varmap[k]=returnString.String()
			    	}else
			    	if strings.Contains(ip,"reserveItem"){
			    		
			    		//z := strings.Split(ip, "(")
				    	//y := strings.Split(z[1], ")") //y[0]
						//w := strings.Split(y[0], ",") //y[0]
						//value :=tagextractor(data, w[0])
						
						varmap[k]=reserveItem(data)//(value,delimiter,index,tailer)
						//fmt.Println("Extract Caled",varmap[k])
					
						
			    	}else{
//			    		fmt.Println("before tag extraction")
			    		varmap[k]=tagextractor(data, ip)
//			    		fmt.Println("value for tag")
//			    		fmt.Println(varmap[k])
			    	}
   	}
  //	fmt.Println(varmap)
  	return varmap
} 

