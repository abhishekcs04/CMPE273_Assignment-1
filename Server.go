package main

import 
(
 	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"	
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"	
)
type meta struct  { 
    Type string
    Start int 
    Count int
}
type Output struct {
	
   Transaction_Number int
   RemainingAmount float64
   CurrentValue float64
   Symbol [20]string
   Shares [20]string
   Price [20]float64
   Quantity [20]int
 }

type resource struct { 

   Classname string
   Fields fields
}
type resources []struct { 
  Resource resource
}

type list struct { 
   Meta meta 
   Resources resources
}	
type fields struct { 
 Name string
 Price string
 Symbol string
 Ts string
 Type string
 Utctime string
 Volume string
}
type jsonType struct { 

   List list
}


// ***************Json format

type Transactions struct { 

  TradeID int
  ReturnBudget float64
  Quantity [200]int
  Price [20]float64
  Shares [200]string
}
type FinancialInfo struct { 
  
   Trade_Number int
}
type Input struct { 
    
   Budget float64 
   Percent []float64
   Shares []string   
}

type Arith struct { }  

var Global_Transaction []Transactions

func FinancialDetails(instrument_code string) []string { 

      url := "http://finance.yahoo.com/webservice/v1/symbols/"+instrument_code+"/quote?format=json"
	  	 
      response_object, err := http.Get(url)
      
      body, err1 := ioutil.ReadAll(response_object.Body)
	  
	  var json_Variable jsonType
      err_var := json.Unmarshal(body, &json_Variable)
	
	  // letters := []string{"a", "b", "c", "d"}
	
       Details_of_Stocks := []string{json_Variable.List.Resources[0].Resource.Fields.Name,
		                              json_Variable.List.Resources[0].Resource.Fields.Symbol,
									  json_Variable.List.Resources[0].Resource.Fields.Price}
											
		 	if err1 != nil || err != nil ||err_var !=nil { 
	
	     fmt.Println(err1)
	     fmt.Println(err)
	     fmt.Println(err_var)
	  
	     response_object.Body.Close()  
	}	 
	    return Details_of_Stocks  
}

var Local_Counter int =0
var Global_Trade_Counter int = 90000

func FinancialBudget(Budget_info float64, Percent_total float64) float64 { 

    information :=  (Budget_info * Percent_total) / 100
	
    return information
}
func PriceofStocks(stock string) string  { 
       Details := FinancialDetails(stock) 
	   //fmt.Printl(Details);
	   return Details[2]   
  }

func Check_Incr_Count(status bool)  {
	
	if status == true { 
	     
		   Global_Trade_Counter++
		   Local_Counter++ 
	    
		}  else {  
	                fmt.Println("Transaction not Incremented")
	            }
	
}

func (t *Arith) Trade(args *Input, out *Output) error {
	
	arguments_Stocks :=len(args.Shares)-1
	
	Budget_value :=0.00
	flag := false
    
	for i:=0 ; i<arguments_Stocks;i++ { 
	  
	    Temp_Current_Price := PriceofStocks(args.Shares[i])
		
		Current_Price,abnormal_err  :=  strconv.ParseFloat(Temp_Current_Price,64)
		 
		if abnormal_err!=nil { 
		
		     fmt.Println("Cannot connect to Yahoo API..")
		     fmt.Println(abnormal_err)  
		}
		
        Temp_Budget := FinancialBudget(args.Budget,args.Percent[i])	
					
		Total_Quantity := Temp_Budget/Current_Price
		
		if Total_Quantity < 1 {  
		
		       fmt.Println("Cannot buy Stock Quantity less than 1")
		
			     
		   }  else  { 
		          	
				   Global_Transaction[Local_Counter].Price[i] = Current_Price
				   Global_Transaction[Local_Counter].TradeID = Global_Trade_Counter
				   Global_Transaction[Local_Counter].Shares[i] = args.Shares[i]
				   Absolute_Quantity := int(Total_Quantity)
				   Temp_Qty := float64(Absolute_Quantity)
				   Global_Transaction[Local_Counter].Quantity[i] = int(Total_Quantity)
				 
				
				   Budget_value = Budget_value + (Temp_Qty * Current_Price)
				   Remaining_budget := args.Budget-Budget_value
		           Global_Transaction[Local_Counter].ReturnBudget = Remaining_budget
			
				   out.Price[i] = Current_Price
				   out.Transaction_Number = Global_Trade_Counter
				   out.Shares[i] = args.Shares[i]
				   out.Quantity[i] = int(Total_Quantity)
				   out.RemainingAmount = Remaining_budget
				   flag = true
				  
			} 
									
    }
    
	
	Check_Incr_Count(flag) 		 
           
				
					
   	return nil
}

func (t *Arith) FinancialInfo(args *FinancialInfo, info *Output) error { 
   
     for var_c:=0;var_c< len(Global_Transaction);var_c++ { 
	       
		   if Global_Transaction[var_c].TradeID == args.Trade_Number  {
		     
				for count:=0; count < len(Global_Transaction[var_c].Shares) ;count++ {
				
				if Global_Transaction[var_c].Shares[count] == "" { 
				     
			        break
				}
				CurrentStockPrice,abnormal_err  :=  strconv.ParseFloat(PriceofStocks(Global_Transaction[var_c].Shares[count]),64)
				info.Price[count] = CurrentStockPrice
				
				if abnormal_err!=nil { 
				
				    fmt.Println("Cound not connect to Yahoo API.. ",abnormal_err)
					
					}
				
				info.Symbol[count] =""
				info.Quantity[count] = Global_Transaction[var_c].Quantity[count]
				info.Shares[count] = Global_Transaction[var_c].Shares[count]
				info.Transaction_Number = Global_Transaction[var_c].TradeID
				
				return_tmp_budget := Global_Transaction[var_c].ReturnBudget 
		        info.RemainingAmount = return_tmp_budget
				
				temp_Budget := float64(Global_Transaction[var_c].Quantity[count])
				
				info.CurrentValue = info.CurrentValue + (temp_Budget * CurrentStockPrice)
				
				if CurrentStockPrice > Global_Transaction[var_c].Price[count] { 
				
				    info.Symbol[count] = "+"
				
				} else if CurrentStockPrice < Global_Transaction[var_c].Price[count] { 
				   
				    info.Symbol[count] = "-"    
					
				}
			  } 
			   
			  
		   }
		   
	   }
	
   return nil
}

func startServer() {
    arith := new(Arith)
	server := rpc.NewServer()
	server.Register(arith)
	rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":1234")
    if e != nil {
        log.Fatal("Error while listening to port -- Error Code :", e)
    }
    for {
        connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error while running the server. Try Again..")
            log.Fatal(err)
        }
    go server.ServeCodec(jsonrpc.NewServerCodec(connection))
    }
}

func main()  {

Global_Transaction=make([]Transactions,50)
go startServer()
fmt.Println("Server ON ...")
var wait string
fmt.Scanln(&wait)

} 