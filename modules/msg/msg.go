/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package msg

type Message struct {
	MsgType    int
	MsgContent interface{}
}

const (
	MsgTypeStop       = 1
	MsgTypeStart      = 2
	MsgTypeLimitSpeed = 3
	MsgQueryStatus    = 4
	MsgRepportStatus  = 5
)
