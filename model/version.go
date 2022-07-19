package model

// ⛵ 🥥 🍺
// 🚢 🚀 🚟
var (
	LOGO = `
                         _ _  
    __ _  ___  ___  __ _(_) | 
   / _  |/ _ \/ _ |/ _  | | | 
  |  g  |  o  \ s \  a  i   l 
   \__, |\___/|___/\__,_|_|__|
   |___/   	
`
	EMOJI    = map[string]string{"gosail": "⛵", "gocy": "🥥", "gobars": "🍺"}
	Historys = []History{
		{Version: "0.1.0",
			Description: "first version",
		},
		{Version: "0.2.0",
			Description: "functional testing",
		},
		{Version: "0.3.0",
			Description: `no-secret sign-in : ues default key path "UserHomeDir/.ssh/id_rsa"`,
		},
		{Version: "0.4.0",
			Description: `enter the id to select the host login`,
		}, {Version: "0.4.1",
			Description: `compatible host by username@host`,
		}, {Version: "0.4.2",
			Description: "non-linux mode show cmdline behind hostname, must show the commands",
		}, {Version: "0.4.3",
			Description: "show exec commands",
		}, {Version: "0.4.4",
			Description: "add a parameter '-s' as a switch for whether or not to use the selected host",
		}, {Version: "0.4.5",
			Description: "can sikp unconnect host",
		},
		{Version: "0.5.0",
			Description: "collated log output",
		}, {Version: "0.5.1",
			Description: "fix ipfiles CIDR'problem",
		},
		{Version: "0.6.0",
			Description: "server list : show the host status",
		}, {Version: "0.6.1",
			Description: "change terminal size , not dynamically",
		}, {Version: "0.6.2",
			Description: "limit ssh connect by WaitGroup",
		}, {Version: "0.6.3",
			Description: "simplify output contenty",
		}, {Version: "0.6.4",
			Description: "adjusted server list",
		}, {Version: "0.6.5",
			Description: "select host : fix selecting host again requires an additional character",
		},
		{Version: "0.7.0",
			Description: "gocy : new scp function , use push/pull copying files in different scenarios",
		}, {Version: "0.7.1",
			Description: "gocy's local files support relative path",
		}, {Version: "0.7.2",
			Description: "fix username' problem ()",
		}, {Version: "0.7.3",
			Description: "fix use relative path when host and user is not nil",
		},
		{Version: "0.8.0",
			Description: "refactored the code logic and operation of gocy",
		}, {Version: "0.8.1",
			Description: "fix gocy pull/push local copy's error; text left-aligned display",
		}, {Version: "0.8.2",
			Description: "modify file permissions 0777",
		},
		{Version: "0.9.0",
			Description: "gobars : k8s mulit exec command tool",
		}, {Version: "0.9.1",
			Description: "add tools' Emoji, add result and selection's username",
		}, {Version: "0.9.2",
			Description: "fix splice command line not executing problem",
		}, {Version: "0.9.3",
			Description: "add gobars'copy function : copy pods' file to local",
		}, {Version: "0.9.4",
			Description: "fix error: username is not specified error",
		}, {Version: "0.9.5",
			Description: "move goscp limit function to client",
		}, {Version: "0.9.6",
			Description: `change code's "sur" to "src", "dst" to "dest"`,
		}, {Version: "0.9.7",
			Description: "gobars' result show the container name",
		}, {Version: "0.9.8",
			Description: "fix gocy root not specificed",
		},
		{Version: "0.10.0",
			Description: "gobars/gocy -copy -tar : gobars and gocy's copy add tar function",
		}, {Version: "0.10.1",
			Description: "reduce unnecessary log output",
		}, {Version: "0.10.2",
			Description: "fix gobars pod name regex pattern",
		}, {Version: "0.10.3",
			Description: "modifily version history",
		}, {Version: "0.10.4",
			Description: "gobars -shell : specifiy container' shell , default sh",
		}, {Version: "0.10.5",
			Description: "gobars -hightlight : highlight output's key",
		}, {Version: "0.10.6",
			Description: "fix  gocy show result without tar problem",
		}, {Version: "0.10.7",
			Description: "fix gobars pull file's name error problem",
		},
		{Version: "0.11.0",
			Description: "new command line interaction with cobra",
		}, {Version: "0.11.1",
			Description: "fix k8s pods entry error problem",
		},
		{Version: "0.12.0",
			Description: "new command line interaction with grumle",
		}, {Version: "0.12.1",
			Description: "Modify the login related parameter names",
		}, {Version: "0.12.2",
			Description: "interactive commands can record the path of the last execution",
		}, {Version: "0.12.3",
			Description: "host exec highlight args support for cmdline and linuxmode",
		}, {Version: "0.12.4",
			Description: "fix gobars container name and gosail -v problem",
		}, {Version: "0.12.5",
			Description: "fix interactive commands pwd get probelm",
		}, {Version: "0.12.6",
			Description: "fix config execution probelm",
		},
	}
)

// furture : dynamic display of the command execution process
