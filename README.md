# TaxenHeimer
As the name says it, this is a project CoppenHeimer clone I made in Golang

The scanner uses IP ranges to scan for Minecraft servers either on the default port (25565)
or it finds the port by looking at the SRV record of the domain.

For the time being, the scanner is not multithreaded and I need help to find a way to avoid deadlocks

## How to run
```bash
git clone https://github.com/TaxMachine/TaxenHeimerGo.git
go get
go build
./TaxenHeimer
```