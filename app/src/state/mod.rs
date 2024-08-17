use axum::extract::FromRef;

use self::traits::TryFromEnv;

pub mod config;
pub mod traits;

#[derive(Clone, FromRef)]
pub struct AppState {
    pub config: std::sync::Arc<config::State>,
}

impl TryFromEnv for AppState {
    fn try_from_env() -> anyhow::Result<Self> {
        Ok(AppState {
            config: config::State::try_from_env()?.into(),
        })
    }
}
