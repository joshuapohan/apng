# apng

## APNG encoder implemented in Go

### How to use
1. Instantiate apng.APNGModel
2. Append images (File)
   ```go 
   APNGModel.AppendImage() 
   ```
   and the corresponding delay for that image 
   ```go 
   APNGModel.AppendDelay()
   ```
3. Run 
   ```go 
   APNGModel.Encode()
   ```
4. Run 
   ```go 
   APNGModel.SaveAsPNG(path) 
   ```
   with path being the target filename (string) to save the apng image
