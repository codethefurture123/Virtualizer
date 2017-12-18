package main

import (
    "fmt"
	"net/http"
	"io/ioutil"
	"strings"
	 "math/rand"
   "time"
    "github.com/Sirupsen/logrus"
//    "github.com/clbanning/mxj"
)

func getNewhandler(service Service) http.HandlerFunc {
	//fmt.Println("inside getHandler")
	return func (w http.ResponseWriter, r *http.Request) {
			//fmt.Println("inside handeler")
			start:=time.Now()
			var operation Operation
			var operationMulti Operation
			var response string
			var evaluatedIPVariables map[string]string
			 w.Header().Set("content-type", service.Type)
			 hah, err := ioutil.ReadAll(r.Body)
	         if err != nil {
		       fmt.Fprintf(w, "%s", err)
	         }
	         reqbody := string(hah)
//	         fmt.Println("-----reqbody:")
//	         	fmt.Println(reqbody)
//)){//||strings.Contains(r.RequestURI, soperation.Opname&& soperation.MultipleResponse==0
here:
	         for _,soperation := range service.Operations{	         	
				if(strings.Contains(reqbody, soperation.Opname)&& soperation.MultipleResponse==0) {
				operation=soperation	
				}else if(strings.Contains(reqbody, soperation.Opname)&&soperation.MultipleResponse!=0){
				 
				// for _,soperation := range service.Operations{
				operationMulti=soperation	
				for _,output := range operationMulti.Outputs {
				if(strings.Contains(output.Tagvalue,"=")){
	          			correlationTagValue:=strings.Split(output.Tagvalue, "=")
	          			xval :=tagextractor(hah, correlationTagValue[0])
	          			//fmt.Println("value found in multiresponse =",xval)
	          			
	          			
	          			//fmt.Println("value to be compared",correlationTagValue[1])
						if(xval==correlationTagValue[1]){
						operation=soperation
						//fmt.Println("Operationnames :",operation.Opname)
						//fmt.Println("before break")
						break here
						
						}
						}else if(strings.Contains(output.Tagvalue,"~")){
	          			correlationTagValue:=strings.Split(output.Tagvalue, "~")
	          			
	          			xval :=tagextractor(hah, correlationTagValue[0])
	          			//fmt.Println("correlationTagValue ~",correlationTagValue)
	          			//fmt.Println("value found in multiresponse ~",xval)
	          			
						if(strings.Contains(xval,correlationTagValue[1])){
						operation=soperation
						break here
						}
						//} 
						} else if (output.Tagvalue=="") {
						operation=soperation
						break here
						}
				}
			
				}
	          }
	          
	         if operation.Monitoring{
				monitoringRequestLogger(start,r.Header,reqbody,operation.Opname,r.RequestURI)
			  }

			  if (operation.Opname!=""){
			  	//fmt.Println("Operationnames :",operation.Opname)
			  	// Request validation(to check whether the request has correct operation name or any essential data)
			  	
			rseed:=time.Now().UnixNano()
			rand.Seed(rseed)
			
	         //output:=operation.Output
	         logrus.WithFields(logrus.Fields{"configured outputs":operation.Outputs}).Debug()
	         
	          for _,output := range operation.Outputs {
	          	logrus.WithFields(logrus.Fields{"Output":output.Response}).Debug()
	          	
	          	correlationValue:=output.Tagvalue
	          	//fmt.Println("Before correlationValue:",correlationValue)
	          	        	
	          	if(strings.Contains(reqbody, correlationValue)){
		          	 response=output.Response
			   }else 
	          	if((strings.Contains(correlationValue,".")) && (!(strings.Contains(correlationValue,"=")))&& (!(strings.Contains(correlationValue,"~")))){
	          			xval :=tagextractor(hah, correlationValue)
	          			//fmt.Println("value found in .",xval)
	          			if xval!=""{
						//fmt.Println("value found")
						response=output.Response
						} else {
							 response="Sorry!,This request is not properly mapped to a response. Please check if the server is configured with a proper Xpath."
						}

	          		}else 
	          	if(strings.Contains(correlationValue,"=")){
	          			correlationTagValue:=strings.Split(correlationValue, "=")
	          			xval :=tagextractor(hah, correlationTagValue[0])
	          			//fmt.Println("value found in =",xval)
						if(xval==correlationTagValue[1]){
						response=output.Response
						} else {
							response="Sorry!,This request is not properly mapped to a response. Please check if the server is configured with a proper Xpath and condition value."
						}
	          	}
				if(strings.Contains(correlationValue,"~")){
	          			correlationTagValue:=strings.Split(correlationValue, "~")
	          			xval :=tagextractor(hah, correlationTagValue[0])
	          			//fmt.Println("value found in ~",xval)

						if(strings.Contains(xval,correlationTagValue[1])){
						response=output.Response
						} else {
							response="Sorry!,This request is not properly mapped to a response. Please check if the server is configured with a proper Xpath and condition value."
						}
						}
	          	evaluatedIPVariables=evaluateInputVariables(output.Variables,hah)
			
	 }
		  
			startdelimiter:= "${"
			enddelimeter:="}"
           // fmt.Println(strings.Contains(output, delimiter))

             for strings.Contains(response, startdelimiter){
             
	         z := strings.SplitN(response, startdelimiter,2)
	         y := strings.SplitN(z[1], enddelimeter,2)
	        // fmt.Println("response splits...",evaluatedIPVariables[y[0]])
	        // fmt.Println("response splits...",y[0])
	         response=z[0]+evaluatedIPVariables[y[0]]+y[1]
	              
   
			}  
//               
//			 <-ch
			 time.Sleep(time.Duration(operation.Delay)* time.Second)
			  //fmt.Println(response)
			//fmt.Fprintf(w, response)
	} else {
		response=`operation not found`
		

	}
	fmt.Fprintf(w, response)
	
	if operation.Monitoring{
		monitoringResponseLogger(time.Now(),w.Header(),response)
	}

		
			        
    }
}