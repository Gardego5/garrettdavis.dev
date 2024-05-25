use std::io::Write;

pub struct Base64Body<T>(pub T);

impl <T: maud::Render> maud::Render for Base64Body<T> {
    fn render_to(&self, buffer: &mut String) {
        let mut enc = base64::write::EncoderStringWriter::from_consumer(buffer, &base64::engine::general_purpose::STANDARD);
        write!(enc, "{}", self.0.render().0).unwrap()
    }
}
