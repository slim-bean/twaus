# The World Around US

## Sensors

### O2

### CO2

### PM 2.5, 10 + VOC

```
git remote add sen54 https://github.com/Sensirion/embedded-i2c-sen5x.git
git subtree add --prefix pkg/sensors/sen54/driver sen54 master --squash
git fetch sen54 master
git subtree pull --prefix pkg/sensors/sen54/driver sen54 master --squash
```

### VOC

### Radiation

### Temp/Humdity

### Barometer