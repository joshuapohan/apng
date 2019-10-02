# apng

## APNG encoder implemented in Go

### How to use
1. Instantiate apng.APNGModel
2. Append images (File)	APNGModel.AppendImage(f) and the corresponding delay for that image APNGModel.AppendDelay(f)
3. Run APNGModel.Encode()
4. Run APNGModel.SaveAsPNG(path) with path being the path to save the apng image
