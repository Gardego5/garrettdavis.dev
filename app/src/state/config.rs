use super::traits::TryFromEnv;

#[derive(Debug, Clone)]
pub struct State {
    pub port: u16,
    pub host: String,
    pub static_dir: String,
    pub data_dir: String,
}

impl TryFromEnv for State {
    fn try_from_env() -> anyhow::Result<Self> {
        use std::env::var;
        let state = State {
            port: var("PORT").unwrap_or("3000".into()).parse::<u16>()?,
            host: var("HOST").unwrap_or("0.0.0.0".into()),
            static_dir: var("STATIC_DIR")?,
            data_dir: var("DATA_DIR")?,
        };
        Ok(state.into())
    }
}
