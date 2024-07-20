resource "aws_cloudfront_distribution" "garrettdavis_dev" {
  aliases = ["garrettdavis.dev", "www.garrettdavis.dev"]

  origin {
    origin_id = "garrettdavis_dev_lb"
    custom_origin_config = {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2", "TLSv1.3"]
    }
  }

  default_cache_behavior {
    allowed_methods        = []
    cached_methods         = []
    target_origin_id       = "garrettdavis_dev_lb"
    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
    }
  }
}
