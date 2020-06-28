systemctl stop apache2
systemctl stop mysql
sync; echo 3 > /proc/sys/vm/drop_caches
swapoff -a && swapon -a
