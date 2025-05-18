# ScyllaDB config

kubectl label nodes uoc-rock3a-02 scylla.scylladb.com/node-type=scylla

## Prepare persistent volumes


radxa@uoc-rock3a-01:~$ sudo cat /etc/fstab
[sudo] password for radxa: 
UUID=502efe26-7c43-49ff-bfce-9e9c230defdb	/	ext4	defaults	0	1
UUID=91C6-5B20	/config	vfat	defaults,x-systemd.automount	0	2
UUID=91C9-60C4	/boot/efi	vfat	defaults,x-systemd.automount	0	2
/dev/sda1 /mnt/persistent-volumes	xfs	auto,nofail,noatime,rw,user,pquota	0	0
#/dev/nvme0n1p1	/mnt/persistent-volumes	xfs	defaults,prjquota	0	1
