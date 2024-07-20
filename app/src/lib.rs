use std::io::Write;

use aws_lambda_events::{apigw::ApiGatewayV2httpResponse, encodings::Body};
use base64::{engine::general_purpose::STANDARD, write::EncoderStringWriter};
use http::{header::CONTENT_TYPE, HeaderMap, HeaderValue};
use maud::Render;

pub mod components;
pub mod pages;

pub struct HtmlResponse(ApiGatewayV2httpResponse);

impl HtmlResponse {
    pub fn new<T: Render>(page: T) -> Self {
        let text_html: HeaderValue = HeaderValue::from_str("text/html").unwrap();

        let mut enc = EncoderStringWriter::new(&STANDARD);
        write!(enc, "{}", page.render().0).unwrap();

        let mut headers = HeaderMap::new();
        headers.insert(CONTENT_TYPE, text_html);

        Self(ApiGatewayV2httpResponse {
            status_code: 200,
            headers,
            multi_value_headers: HeaderMap::new(),
            body: Some(Body::Text(enc.into_inner())),
            is_base64_encoded: true,
            cookies: vec![],
        })
    }

    pub fn headers(&self) -> &HeaderMap {
        &self.0.headers
    }
    pub fn headers_mut(&mut self) -> &mut HeaderMap {
        &mut self.0.headers
    }

    pub fn multi_value_headers(&self) -> &HeaderMap {
        &self.0.multi_value_headers
    }
    pub fn multi_value_headers_mut(&mut self) -> &mut HeaderMap {
        &mut self.0.multi_value_headers
    }

    pub fn cookies(&self) -> &Vec<String> {
        &self.0.cookies
    }
    pub fn cookies_mut(&mut self) -> &mut Vec<String> {
        &mut self.0.cookies
    }

    pub fn status(&self) -> &i64 {
        &self.0.status_code
    }
    pub fn status_mut(&mut self) -> &mut i64 {
        &mut self.0.status_code
    }
}

impl From<HtmlResponse> for ApiGatewayV2httpResponse {
    fn from(value: HtmlResponse) -> Self {
        value.0
    }
}
