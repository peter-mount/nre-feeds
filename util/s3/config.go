package s3

type S3Credentials struct {
  // The bucket name
  Bucket        string    `yaml:"bucket"`
  // Prefix with trailing / to prepend to all file names within the bucket
  Path          string    `yaml:"path"`
  // Access Key
  AccessKey     string    `yaml:"accessKey"`
  // Secret Key
  SecretKey     string    `yaml:"secretKey"`
  // Bucket region
  Region        string    `yaml:"region"`
}
