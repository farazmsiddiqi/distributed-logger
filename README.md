# CS425 MP0

SEE FULL REPORT W/ GRAPHS HERE: https://docs.google.com/document/d/1CS3yDMqGTkC-9xEDM5SMQ1CNyVjrYCNZMPJv0UZKaP8/edit?usp=sharing 

Members: Mahi Kolla (mkolla2), Faraz Siddiqi (farazms2)
Cluster Number: 31
Repository + version number: https://gitlab.engr.illinois.edu/mkolla2/cs425-mp0/-/tree/main 

# Instructions for building and running your code:
on the central logger VM (sp23-cs425-3101.cs.illinois.edu) run 
go run central_logger.go [port_number]
on each of your node VMs (sp23-cs425-3102.cs.illinois.edu onward) run 
python3 -u generator.py [hz] | go run node.go [node_name] [address of vm] [same port_number as before!!]
use CTRL+C to close each process 

# Description of how you are measuring the delay and bandwidth:
Delay: We measure delay of an event being logged on the central logger by subtracting the time the event was generated (which is sent to us from the node) from the current time at the central logger. This delay is calculated and then written to an aux_log file. 

event delay = [current time at central logger] - [time event was generated as sent by node]

Bandwidth: We measure bandwidth used by the central logger by recording the length of each message it receives from the nodes. 
