# dyndns
application to update dns records at digital ocean

success/failure/no action recorded to /var/log/syslog

## expects a json file  ~/.config/dyndns
```
{
  Token: "digitaloceantoken",
  Host: "hostname",
  Id: "digitalocean dns record identifier"
}
```
## usage
place file in /home/user/bin/

create cron entry
```
# run crontab -e
# sample cron job task 

# m h  dom mon dow   command
*/5 * * * * /home/user/bin/ip
```
