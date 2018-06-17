package statsdclient


import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
	"strings"
	"log"
)


type Client struct {
	Address string
	Timeout int
	Prefix string
	Connection io.WriteCloser
}

func NewClient(address string) *Client {
	return &Client{Address: address}
}

func (c *Client) Connect() (err error) {
	if c == nil {
		return nil
	}

	if c.Connection != nil {
		return fmt.Errorf("client is already connected. Call Disconnect first if you want to reconnect")
	}




	c.Connection, err = net.DialTimeout("udp", c.Address, 40 * time.Second)
	return err
}

// Close implements the io.Closer interface by closing the clients Connection.
// If c is nil this will do nothing (noop)
func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	err := c.Connection.Close()
	c.Connection = nil
	return err
}


func (c *Client) Send(metric Metric)  {

		select {
		case  <-time.After(5 * time.Second):
			log.Printf("timeout sending metric to stastsd for host %")
			return
		default:
			if c == nil {
				fmt.Errorf("well... client is broken. Try to make an instance of it and call a Connect method.")
				return
			}

			buffer := &bytes.Buffer{}
			Name := fmt.Sprintf("%s.%s", c.Prefix, metric.Name)
			Name= strings.Replace(Name, ".", "_", -1)
			_, err:= fmt.Fprintf(buffer, "%s:%v|g", Name, metric.Value)
	//		fmt.Printf("%s : %v |g\n", Name, metric.Value)
			if err != nil {
				log.Printf("Error write to buffer %s", err.Error())
				return
			}
			_, err = c.Connection.Write(buffer.Bytes())
			if err != nil {
				log.Printf( "Error write to conn %s",err.Error())
				return
			}
		}

}

