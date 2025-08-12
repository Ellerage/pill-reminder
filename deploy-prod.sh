# VPS
ssh $VPS_SSH << 'PREP'
sudo systemctl stop pill-reminder
rm ./pill-reminder
rm ./pill-reminder.log
redis-cli <<REDISCLEANUP
SELECT 0
FLUSHDB
SELECT 1
FLUSHDB
EXIT
REDISCLEANUP
PREP

# My PC
rm ./pill-reminder &&
make build_prod
scp ./pill-reminder $VPS_SSH:/root/


# VPS
ssh $VPS_SSH "sudo systemctl start pill-reminder"