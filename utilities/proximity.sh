#!/bin/bash
#set -o verbose sh -v
# Copied from Steven on http://gentoo-wiki.com/Talk:TIP_Bluetooth_Proximity_Monitor

# These are the sections you'll need to edit

 # You'll need to use the MAC address of your phone here
DEVICE="00:0F:46:FF:05:07" 

# How often to check the distance between phone and computer in seconds
CHECK_INTERVAL=2

# The RSSI threshold at which a phone is considered far or near
THRESHOLD=-13

# The command to run when your phone gets too far away
FAR_CMD='/opt/gnome/bin/gnome-screensaver-command --activate'

# The command to run when your phone is close again
NEAR_CMD='/opt/gnome/bin/gnome-screensaver-command --poke'

HCITOOL="/usr/bin/hcitool"
STARTX_PID=0
DEBUG="/tmp/btproximity.log"

connected=0

function msg {
    echo "$1" #>> "$DEBUG"
}

function check_connection {
    connected=0;
    found=0
    for s in `$HCITOOL con`; do
        if [[ "$s" == "$DEVICE" ]]; then
            found=1;
        fi
    done
    if [[ $found == 1 ]]; then
        connected=1;
    else
       msg 'Attempting connection...'
        if [ -z "`$HCITOOL cc $DEVICE 2>&1`" ]; then
            msg 'Connected.'
            connected=1;
        else
                if [ -z "`l2ping -c 2 $DEVICE 2>&1`" ]; then
                        if [ -z "`$HCITOOL cc $DEVICE 2>&1`" ]; then
                            msg 'Connected.'
                            connected=1;
                        else
                        msg "ERROR: Could not connect to device $DEVICE."
                        connected=0;
                        fi
                fi
        fi
    fi
}

check_connection

while [[ $connected -eq 0 ]]; do
    check_connection
    sleep 3
done

name=`$HCITOOL name $DEVICE`
msg "Monitoring proximity of \"$name\" [$DEVICE]";

state="near"
while /bin/true; do

    check_connection

    if [[ $connected -eq 1 ]]; then
        rssi=$($HCITOOL rssi $DEVICE | sed -e 's/RSSI return value: //g')

        if [[ $rssi -le $THRESHOLD ]]; then
            if [[ "$state" == "near" ]]; then
                msg "*** Device \"$name\" [$DEVICE] has left proximity"
                state="far"
                $FAR_CMD > /dev/null 2>&1
            fi
        else
            if [[ "$state" == "far" && $rssi -ge $[$THRESHOLD+2] ]]; then
                msg "*** Device \"$name\" [$DEVICE] is within proximity"
                state="near"
                $NEAR_CMD > /dev/null 2>&1
                STARTX_PID=$(pgrep startx)
            fi
        fi
        msg "state = $state, RSSI = $rssi"
    fi

    sleep $CHECK_INTERVAL
done
