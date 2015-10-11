package main

import
(
    "fmt"	
	"net"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"strconv"
	"net/rpc"
)

type Reply struct {
    
	Transaction_Number int
    Shares []string
	Price []float64
    Quantity []int
    RemainingAmount float64
	Prft_Loss [5]string
	CurrentValue float64
}

type TradeDetails struct { 
     Trade_Number int64	
}

type Trade struct { 

  Shares[] string
  Quantity[] string 
  
}

type Transaction_Customers struct { 
 
   Shares[] string 
   Percent[] float64
  
   Budget float64
}
    var trans_info *Transaction_Customers
	var transaction_id *TradeDetails
	var sum float64
	var reply Reply
    var connection *rpc.Client
	var Splited_info []string
    
/* func Initialize_Connection() (*rpc.Client) { 

conn, err := net.Dial("tcp", "localhost:1234")

    if err != nil {
        panic(err)
    }
    defer conn.Close()

    connection := jsonrpc.NewClient(conn)
	return connection

} */

func Transaction_Execution(All_Stocks []string,All_Percentage []float64,Total_Budget float64) { 

// conecc := Initialize_Connection()

conn, err := net.Dial("tcp", "localhost:1234")

    if err != nil {
          fmt.Println("Error while connection to remote server. Please check if Server.go started properly")  
   
		panic(err)
    }
    defer conn.Close()
    conecc := jsonrpc.NewClient(conn)

// var err error
    trans_info = &Transaction_Customers{All_Stocks,All_Percentage,Total_Budget}
	err = conecc.Call("Arith.Trade",trans_info,&reply)
   	
    if err != nil { 
       fmt.Println("Error while connection to remote server. Please check if Server.go started properly")  
    } 
	
	fmt.Println("Transaction ID :",reply.Transaction_Number)
	fmt.Println()
	
	fmt.Print("Portfolio :   ")
	total_Shares := reply.Shares
	
	for x_count:=0;x_count<len(total_Shares);x_count++ {
	   
	   if reply.Shares[x_count] == "" { 
	     
		  break 
	   
	   }
      
	  fmt.Print("[\"",reply.Shares[x_count])
	  fmt.Print(":",reply.Quantity[x_count])
	  fmt.Print(":$",reply.Price[x_count],"\"]")
	
	}
	
	fmt.Println()
	fmt.Println()
	
	fmt.Println("Remaining Budget : $",reply.RemainingAmount) 
	

}
	
func Transaction_Details(Transaction_Number int64) { 

conn, err := net.Dial("tcp", "localhost:1234")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    conecc := jsonrpc.NewClient(conn)

	transaction_id = &TradeDetails{Transaction_Number}
	
	err = conecc.Call("Arith.FinancialInfo",transaction_id,&reply)
	
	fmt.Println("Transaction Number :",reply.Transaction_Number)
	fmt.Println()
	fmt.Print("Portfolio (Stock Holdings) :")
	
	for count:=0;count<len(reply.Shares);count++ {
	   
	   if reply.Shares[count] == "" { 
	     
		  break 
	    
		}
     
	    fmt.Print(reply.Shares[count])
	    fmt.Print(":",reply.Quantity[count])
		fmt.Print()
	    fmt.Print(":",reply.Prft_Loss[count],"$",reply.Price[count],"  ")

	}
	    fmt.Println()
	
	    fmt.Println()
	    fmt.Println("Market Capitilization : ",reply.CurrentValue)
	    fmt.Println()

	    fmt.Println("Remaining Budget : $",reply.RemainingAmount)
	
}

func main() {
   
	// Command Line Arguments
	
	newst := strings.Replace(os.Args[1], "%", "", -1)
	var Parsed_Info []string = strings.Split(newst,",")
    
	Passed_Percentage  := make([]float64,len(Parsed_Info))
    
	Get_percent := Parsed_Info[len(Parsed_Info)-1]
	
	Total_valuation ,err := strconv.ParseFloat(Get_percent,64)
	
    if err!=nil { 
	          fmt.Println("Cannot convert from Int to Float")
	          fmt.Println(err)
	} 
	
	Collected_Stocks := make([]string,len(Parsed_Info))
	
	for i:=0; i<len(Parsed_Info); i++  {
	 
        Splited_info = strings.Split(Parsed_Info[i],":")
		   
		   for j:=0; j<len(Splited_info)-1; j++ { 
		       
		        Collected_Stocks[i] = Splited_info[0]
				
			    Val_Percent := Splited_info[1]
			
			 if temp_Float, err := strconv.ParseFloat(Val_Percent, 64); err == nil {
				       
						Passed_Percentage[i] = temp_Float
						sum = sum + Passed_Percentage[i]
			           
	          }
			  
		  }   
		 
		
	 }
	
	 if(sum != 100 && strings.ContainsAny(os.Args[1],"%" )){
	 fmt.Println("The complete Percent should be 100")
		return
	 }

	if len(Parsed_Info) < 2 { 

    
    Int_Transaction_Number ,err_conver_int:= strconv.ParseInt(Parsed_Info[0],10,64)
	
	if err_conver_int!= nil {
	   fmt.Println("Please enter TRANSACTIION Number and not Float ID.")   
	   fmt.Println(err_conver_int)  
	}
	
	Transaction_Details(Int_Transaction_Number)
	
	} else { 

      Transaction_Execution(Collected_Stocks,Passed_Percentage,Total_valuation)

	
		}	
}