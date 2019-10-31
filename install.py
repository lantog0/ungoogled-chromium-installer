#!/usr/bin/python3

import os
import urllib.request

from sys import exit
from re import search, findall
from threading import Thread

def create_download_dir():
    if not os.path.exists(download_dir):
        os.makedirs(download_dir)

def download_package(package_url):
    package_name = package_url.split('/')[-1]

    print('[-] Downloading {}...'.format(package_name))

    package_bin = urllib.request.urlopen(package_url).read()

    with open(download_dir + package_name, 'wb') as f:
        f.write(package_bin)

    print('[+] File downloaded {}'.format(package_name))

    return

def fetch_packages(version):
    packages_url = download_url + version
    regex = r'https://github.com/Eloston/ungoogled-chromium-binaries/releases/download/[^\"]*?\.deb'

    html_page = urllib.request.urlopen(packages_url).read().decode('utf-8')
    packages = findall(regex, html_page)

    return packages

if __name__ == '__main__':
    installed_version = os.popen('apt search ungoogled-chromium 2>/dev/null | grep -E "[0-9\.-]+.buster1" -o | head -n 1').read().strip()
    download_url = 'https://ungoogled-software.github.io/ungoogled-chromium-binaries/releases/debian/buster_amd64/'

    html_page = urllib.request.urlopen(download_url).read().decode('utf-8')

    try:
        current_version = search('[\d\.-]+.buster1', html_page).group()
    except AttributeError:
        print('Error matching version')
        exit(1)

    if current_version == installed_version:
        packages = fetch_packages(current_version)


    cwd = os.path.dirname(os.path.realpath(__file__))
    download_dir = cwd + '/binaries/'
    create_download_dir()

    threads = []

    for package_url in packages:
        t = Thread(target=download_package, args=(package_url,))
        threads.append(t)
        t.start()

    for t in threads:
        t.join()
