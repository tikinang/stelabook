#!/bin/sh

lscpu >> /mnt/ext1/stelabook.log 2>&1
uname -m >> /mnt/ext1/stelabook.log 2>&1

chmod +x /mnt/ext1/stelabook_client
/mnt/ext1/stelabook_client >> /mnt/ext1/stelabook.log 2>&1
