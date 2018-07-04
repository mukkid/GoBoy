import "strings"
import "os/signal"


type FunctionFrame struct {
    addr uint16 /* Address of frame on stack */
    returnAddr uint16
}

type DebugState {
    breakpoints []uint16
}

func debugLoop() {
    reader := bufio.NewReader(os.Stdin)
    for {
        /* TODO: handle error */
        cmd, _ := reader.ReadString('\n')
        tokens := strings.Fields(cmd)
        if len(tokens) == 0 {
            continue
        }
        switch tokens[0] {
        /* set breakpoint */
        case "b":
            fallthrough
        case "break":
        /* step */
        case "n":
            fallthrough
        case "next":
        /* reset and run */
        case "r":
            fallthrough
        case "run":
        /* continue */
        case "c":
            fallthrough
        case "continue":
        }
    }
}
