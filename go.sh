#!/bin/sh

python3 install.py

cd binaries

sudo dpkg -i ungoogled-chromium_*.deb ungoogled-chromium-common_*.deb
sudo dpkg -i ungoogled-chromium*.deb

cd ..
