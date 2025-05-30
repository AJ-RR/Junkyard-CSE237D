# Building the android source kernel

Reference: 

1. https://source.android.com/docs/setup/build/building-pixel-kernels
2. Gabe’s blog to build pmos - [Porting postmarketOS to the Pixel 6A – Gabriel's Academic and Engineering Blog](https://gabriel.marcanobrady.family/blog/2023/04/24/porting-postmarketos-to-the-pixel-6a/)

aosp branch used by the phones currently: **android-gs-felix-5.10-android14-qpr3**

To download the source

```bash
repo init -u https://android.googlesource.com/kernel/manifest -b **android-gs-felix-5.10-android14-qpr3
repo sync**
```

To create a new kernel config file

```bash
cd aosp
CC=clang ARCH=arm64 LLVM=1 make gki_defconfig
```

To edit kernel config

```bash
CC=clang ARCH=arm64 LLVM=1 make nconfig
```

To save the config

```bash
CC=clang ARCH=arm64 LLVM=1 make savedefconfig
```

The move the config to the correct location

```bash
CC=clang ARCH=arm64 LLVM=1 make mrproper
```

To build, go to the parent folder

```bash
ENABLE_STRICT_KMI=0 BUILD_AOSP_KERNEL=1 ./build_felix.sh
```

Output files will be in “out/mixed/dist”

**To directly modify gki_defconfig - remove check in private/gs-google/build.config.gki**