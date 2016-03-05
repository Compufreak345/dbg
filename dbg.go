// Used for debug output. Switch to false if we release.
// This could also be done by a makefile using different source codes
// (see https://groups.google.com/d/msg/golang-nuts/gU7oQGoCkmg/xlVJx-OJ9EUJ )
package dbg

import (
	"fmt"
	"log"
	"net/http"
	godbg "runtime/debug"
	"strconv"
	"time"
)

// Console formatting inc :

// normal
const KNRM = "\x1B[0m"

// Faint (decreased intensity, not widely supported)
const KFNT = "\x1B[2m"

// Other colors
const KRED = "\x1B[31m"
const KGRN = "\x1B[32m"
const KYEL = "\x1B[33m"
const KBLU = "\x1B[34m"
const KMAG = "\x1B[35m"
const KCYN = "\x1B[36m"
const KWHT = "\x1B[37m"
const KRESET = "\033[0m"

/* const Debugging can be set in production without any security risks, it enables debugging log-output.
 */
const Debugging debug = true

/* const Develop SHOULD NEVER EVER be set in Prod, it is used to disable security-relevant features in the using library
 */
const Develop = true
const v verbose = true
const i info = true
const l lg = true

type debug bool
type verbose bool
type info bool
type lg bool

type Tag string

// http://stackoverflow.com/a/25458067

/*
ERROR – something terribly wrong had happened, that must be investigated immediately. No system can tolerate items logged on this level. Example: NPE, database unavailable, mission critical use case cannot be continued.

WARN – the process might be continued, but take extra caution. Example: “Application running in development mode” or “Administration console is not secured with a password”. The application can tolerate warning messages, but they should always be justified and examined.

INFO – Important business process has finished. In ideal world, administrator or advanced user should be able to understand INFO messages and quickly find out what the application is doing. For example if an application is all about booking airplane tickets, there should be only one INFO statement per each ticket saying “[Who] booked ticket from [Where] to [Where]“. Other definition of INFO message: each action that changes the state of the application significantly (database update, external system request).

DEBUG – Developers stuff.

VERBOSE – Very detailed information, intended only for development. You might keep trace messages for a short period of time after deployment on production environment, but treat these log statements as temporary, that should or might be turned-off eventually. The distinction between DEBUG and VERBOSE is the most difficult, but if you put logging statement and remove it after the feature has been developed and tested, it should probably be on VERBOSE level.
*/

// func D prints a debug message - development only
func D(tag Tag, format string, args ...interface{}) {
	Debugging.P(KNRM+"DEBUG", tag, format, args...)
}

// func V prints a verbose message - development only (importance below debug, e.g. for big variable prints)
func V(tag Tag, format string, args ...interface{}) {
	v.P(KFNT+"VERBOSE", tag, format, args...)
}

// func I prints a info message - will be in production mode, e.g. for registration finished, map uploaded etc.
func I(tag Tag, format string, args ...interface{}) int64 {
	timeKey := time.Now().UnixNano()
	i.P(KGRN+"INFO", tag, strconv.FormatInt(timeKey, 10)+" -- "+format, args...)
	return timeKey
}

// func W prints a warning - something went bad, but the process can be continued. Only allowed in special cases
func W(tag Tag, format string, args ...interface{}) int64 {
	timeKey := time.Now().UnixNano()
	l.P(KYEL+"WARN", tag, strconv.FormatInt(timeKey, 10)+" -- "+format, args...)
	return timeKey
}

// func E prints an error, with stacktrace - IMMEDIATELY FIX THIS!
func E(tag Tag, format string, args ...interface{}) int64 {
	timeKey := time.Now().UnixNano()
	format = fmt.Sprintf(format, args...)
	format += fmt.Sprintf("\n StackTrace : %v", string(godbg.Stack()))
	l.P(KRED+"ERROR", tag, strconv.FormatInt(timeKey, 10)+" -- "+format)
	return timeKey
}

// func WTF prints a WTF - "What a terrible failure"
func WTF(tag Tag, format string, args ...interface{}) int64 {
	timeKey := time.Now().UnixNano()
	l.P(KMAG+"WTF", tag, strconv.FormatInt(timeKey, 10)+" -- "+format, args...)
	return timeKey
}

// func P is used to add a prepositon & tag to the given logmessage & reset colors
func (d debug) P(preposition string, tag Tag, format string, args ...interface{}) {
	if d {
		log.Printf(preposition+"/"+string(tag)+" : "+format+KRESET, args...)
	}
}

// func GetRequest prints a http-Request, if we are in develop the full request, otherwise only Method & url without parameters
func GetRequest(r *http.Request) interface{} {
	if Develop {
		return r
	} else {
		return fmt.Sprintf("[%s] %q %v\n", r.Method, r.URL.Path)
	}
}

// func P is used to add a prepositon & tag to the given logmessage & reset colors
func (l lg) P(preposition string, tag Tag, format string, args ...interface{}) {
	if l {
		log.Printf(preposition+"/"+string(tag)+" : "+format+KRESET, args...)
	}
}

// func P is used to add a prepositon & tag to the given logmessage & reset colors
func (i info) P(preposition string, tag Tag, format string, args ...interface{}) {

	if i {
		log.Printf(preposition+"/"+string(tag)+" : "+format+KRESET, args...)
	}
}

// func P is used to add a prepositon & tag to the given logmessage & reset colors
func (v verbose) P(preposition string, tag Tag, format string, args ...interface{}) {

	if v {
		log.Printf(preposition+"/"+string(tag)+" : "+format+KRESET, args...)
	}
}
