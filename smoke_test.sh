./netcalc --help
./netcalc --manual

# CIDR notation
./netcalc -c 192.168.0.1/24

# IP + mask
./netcalc -i 10.0.0.5 255.255.255.0

# Wildcard mask
./netcalc -w 192.168.0.1/24

# Wildcard → CIDR
./netcalc -W 0.0.0.255

# IP inclusion
./netcalc -I 192.168.0.0/24 192.168.0.5

# List hosts (cuidado: pode gerar muita saída)
./netcalc -l 192.168.0.0/30

# Next/Previous subnet
./netcalc -n 192.168.1.0/24

# Binary display
./netcalc -b 192.168.0.1/24

# Reverse DNS zone
./netcalc -r 192.168.0.0/24
