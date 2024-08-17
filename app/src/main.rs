use axum::{
    routing::{get, post},
    Router,
};
use state::traits::TryFromEnv;
use tower_http::{compression::CompressionLayer, services::ServeDir};

mod components;
mod pages;
mod state;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let state = state::AppState::try_from_env()?;

    let compression_layer = CompressionLayer::new()
        .br(true)
        .deflate(true)
        .gzip(true)
        .zstd(true);

    let listener =
        tokio::net::TcpListener::bind((state.config.host.clone(), state.config.port)).await?;

    let app = Router::new()
        .route("/", get(pages::index::handler))
        .route("/blog", get(pages::blog::handler))
        .route("/contact", get(pages::contact::handle_get))
        .route("/contact", post(pages::contact::handle_post))
        .route("/resume", get(pages::resume::handler))
        .route("/notes", get(pages::notes::handler))
        .route(
            "/presentation/:presentation",
            get(pages::presentations::handler),
        )
        .nest_service("/static", ServeDir::new(state.config.static_dir.clone()))
        .layer(compression_layer)
        .with_state(state);

    Ok(axum::serve(listener, app).await?)
}
