package model

// version history:
// 0.1.0 : first version
// 0.2.0 : functional testing
// 0.3.0 : no-secret sign-in : ues default key path "UserHomeDir/.ssh/id_rsa"
// 0.3.1 : fix ipfile can not use no-secert sign-in
// 0.3.2 : add ssh/terminal.go
// 0.4.0 : enter the id to select the host login
// 0.4.1 : compatible host by username@host
// 0.4.2 : non-linux mode show cmdline behind hostname, must show the commands
// 0.4.3 : show exec commands
// 0.4.4 : add a parameter '-s' as a switch for whether or not to use the selected host
// 0.4.5 : can sikp unconnect host
// 0.5.0 : collated log output
// 0.5.1 : fix ipfiles CIDR'problem
// 0.6.1 : server list : show the host status
// 0.6.2 : change terminal size , not dynamically
// 0.6.3 : limit ssh connect by WaitGroup
// 0.6.4 : Simplify output content

// todo : select host : fix selecting host again requires an additional character
// furture : dynamic display of the command execution process
const (
	VERSION = "0.6.4"
)