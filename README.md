# 🏎️ Gophormula

Currently under active development, Gophormula is a 

## Roadmap

 - [ ] Locate and parse historic data streams
 - [ ] Create a SignalR server to stream historic data
 - [ ] Create a SignalR client to receive historic data
 - [ ] Locate and parse the race calendar to determine the next race
 - [ ] Test the SignalR client on the F1 live timing stream
 - [ ] Expose Prometheus metrics for a race based on a stream
 - [ ] Display race data to users
   - [ ] Grafana Dashboard for Prometheus
   - [ ] SSE Driven website
   - [ ] TUI

## Helpful

 - [FastF1](https://docs.fastf1.dev/) A Python library for parsing live and historic race data.
 - [F1Gopher Lib](https://github.com/f1gopher/f1gopherlib) is a Go library for parsing live and historic race data.
 - [OpenF1](https://openf1.org/) is service that exposes an API for race data in JSON/CSV format. Historical data is free, live is paid. 
 - [Signal R on the Wire](https://blog.3d-logic.com/2015/03/29/signalr-on-the-wire-an-informal-description-of-the-signalr-protocol/)

## Examples

 - [Car Data](https://livetiming.formula1.com/static/2021/2021-04-18_Emilia_Romagna_Grand_Prix/2021-04-18_Race/CarData.z.jsonStream)
 - [Car Data Decompression](https://github.com/theOehrly/Fast-F1/issues/24)
 - [ASP Net SignalR Specification](https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md)