# Connecting to the cluster

# Connecting to automaton:

To get access to automaton, our lab server

- Send your RSA pubkey and desired username to Gabe Marcano <[gmarcano@ucsd.edu](mailto:gmarcano@ucsd.edu)>
- You can then add these lines to your ssh config:
    
    ```bash
    Host automaton
    
    Hostname smartcycling.sysnet.ucsd.edu
    
    port 44422
    ```
    

# Discovering the phones

Phone specs - Google Pixel Fold

Run the following command to discover all phones:

```bash
nmap 10.42.0.0/24
```

The phones have been assigned static IPs in the rage 10.42.0.2-17. The host names are “google-felix-1” to “google-felix-16”. The last part of IP is the host name + 1. For example, the phone “google-felix-3” will have IP 10.42.0.4

The username of all phones is “user” and password is 0000

Note that hostname must be reconfigured if the phone is erased and flashed.

Connect to a phone

```bash
ssh user@10.42.0.<#phone>
```

# Troubleshooting

If no phones are discoverable, it may be due to network card hanging.

```bash
ethtool -K internal0 gso off gro off tso off tx off rx off
```