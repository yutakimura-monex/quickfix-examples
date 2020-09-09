package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
	"github.com/quickfixgo/quickfix"
)

//TradeClient implements the quickfix.Application interface
type TradeClient struct {
}

//OnCreate implemented as part of Application interface
func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

//OnLogon implemented as part of Application interface
func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {
	return
}

//OnLogout implemented as part of Application interface
func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

//FromAdmin implemented as part of Application interface
func (e TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

//ToAdmin implemented as part of Application interface
func (e TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	return
}

//ToApp implemented as part of Application interface
func (e TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("\n")
	fmt.Printf("FromApp: %s\n", msg.String())
	fmt.Printf("as\n")

	msgType, _ := msg.MsgType()
	fmt.Printf("  MsgType: %s \n", msgType)

	tags := msg.Body.Tags()
	sort.Sort(ByTag(tags))
	for _, tag := range tags {
		if val, err := msg.Body.GetString(tag); err == nil {
			fmt.Printf("  %3d: %s\n", tag, val)
		} else if val, err := msg.Body.GetInt(tag); err == nil {
			fmt.Printf("  %3d: %d\n", tag, val)
		} else if val, err := msg.Body.GetTime(tag); err == nil {
			fmt.Printf("  %3d: %s\n", tag, val.String())
		}
	}
	return
}

// orde by Tag
type ByTag []quickfix.Tag

func (t ByTag) Len() int           { return len(t) }
func (t ByTag) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTag) Less(i, j int) bool { return t[i] < t[j] }

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "tradeclient.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	app := TradeClient{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		fmt.Println("Error creating file log factory,", err)
		return
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	initiator.Start()

Loop:
	for {
		action, err := internal.QueryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = internal.QueryEnterOrder(appSettings)

		case "2":
			err = internal.QueryCancelOrder(appSettings)

		case "3":
			err = internal.QueryMarketDataRequest(appSettings)

		case "4":
			//quit
			break Loop

		default:
			err = fmt.Errorf("unknown action: '%v'", action)
		}

		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	initiator.Stop()
}
