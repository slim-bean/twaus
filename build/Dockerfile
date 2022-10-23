FROM golang:1.19-bullseye

RUN apt-get update && apt-get install -y git && apt-get install -y build-essential


# I was unable to use the arm toolchain in bullseye to build for the pi zero, it seems support for arm6 may be missing
# raspberry pi has some really old cross compilers available but it's a very old GCC
# found this project which has a newer GCC https://github.com/Pro/raspi-toolchain
# Usage: CC=/opt/cross-pi-gcc/bin/arm-linux-gnueabihf-gcc
RUN wget https://github.com/Pro/raspi-toolchain/releases/latest/download/raspi-toolchain.tar.gz && tar xfz raspi-toolchain.tar.gz --strip-components=1 -C /opt

