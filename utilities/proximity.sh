#!/bin/bash
#set -o verbose sh -v
# Copied from Steven on http://gentoo-wiki.com/Talk:TIP_Bluetooth_Proximity_Monitor
# ...then stolen from http://www.novell.com/coolsolutions/feature/18684.html for huego

# Usage: proximity.sh <mac address>

 # You'll need to use the MAC address of your phone here
DEVICE="${1:-no MAC set}"

# How often to check the distance between phone and computer in seconds
NEAR_CHECK_INTERVAL=5
FAR_CHECK_INTERVAL=5

# The command to run when your phone gets too far away
FAR_CMD='(echo far; date) >> /tmp/log'

# The command to run when your phone is close again
NEAR_CMD='(echo near; date) >> /tmp/log'

HCITOOL="/usr/bin/hcitool"
DEBUG="/tmp/btproximity.log"

function msg {
    echo "$1" #>> "$DEBUG"
}

name=`$HCITOOL name $DEVICE`
msg "Monitoring proximity of \"$name\" [$DEVICE]";

state="far"
while /bin/true; do
        if l2ping -c 1 ${DEVICE}; then
            if [[ "$state" == "far" ]]; then
                msg "*** Device \"$name\" [$DEVICE] is within proximity"
                state="near"
                $NEAR_CMD > /dev/null 2>&1
            fi
            sleep ${NEAR_CHECK_INTERVAL}
        else
            if [[ "$state" == "near" ]]; then
                msg "*** Device \"$name\" [$DEVICE] has left proximity"
                state="far"
                $FAR_CMD > /dev/null 2>&1
            fi
            sleep ${FAR_CHECK_INTERVAL}
        fi
done
