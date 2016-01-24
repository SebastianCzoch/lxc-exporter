# LXC exporter
LXC exporter is small application written in go which are providing some metrics about LXC containers running on host in Prometheus.io format.
It's beta version, already tested on Ubuntu Willy (15.10) and linux kernel 4.x.x

## Metrics
| Metric name           			| Description                                             					| Enabled by default |
|-----------------------------------|---------------------------------------------------------------------------|--------------------|
| lxc_cpu               			| Seconds the cpus spent in each mode. For all containers 					| yes                |
| lxc_cpu_precentage    			| Precentage of usage processor                           					| yes                |
| lxc_cpu_physical_real 			| Seconds the real physical cpu spent in each mode. (minus containers usage)| yes                |
| lxc_cpu_physical_real_precentage	| Precentage of usage processor (minus containers usage)       				| yes                |
| lxc_memory_usage					| Memory usage in each container in bytes       							| yes                |

## Flags
| Name               	| Description                                 	| Default value 	|
|--------------------	|---------------------------------------------	|---------------	|
| web.listen-address 	| The address to listen on for HTTP requests. 	| :9125         	|

## Building and running

    go build
    ./lxc_exporter <flags>

## Running tests

    go test ./...
