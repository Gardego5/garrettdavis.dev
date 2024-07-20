use anyhow::Result;
use axum::{
    routing::{get, post},
    Router,
};
use garrettdavis_dev::pages;
use tower_http::{compression::CompressionLayer, services::ServeDir};

const VAR_PORT: &'static str = "PORT";
const VAR_HOST: &'static str = "HOST";
const VAR_STATIC_DIR: &'static str = "STATIC_DIR";

#[tokio::main]
async fn main() -> Result<()> {
    let compression_layer = CompressionLayer::new()
        .br(true)
        .deflate(true)
        .gzip(true)
        .zstd(true);

    let app = Router::new()
        .route("/", get(pages::index::handler))
        .route("/blog", get(pages::blog::handler))
        .route("/contact", get(pages::contact::handle_get))
        .route("/contact", post(pages::contact::handle_post))
        .route("/resume", get(pages::resume::handler))
        .route("/notes", get(pages::notes::handler))
        .nest_service("/static", ServeDir::new(std::env::var(VAR_STATIC_DIR)?))
        .layer(compression_layer);

    let listener = tokio::net::TcpListener::bind((
        std::env::var(VAR_HOST).unwrap_or("0.0.0.0".into()),
        std::env::var(VAR_PORT)
            .unwrap_or("3000".into())
            .parse::<u16>()?,
    ))
    .await?;

    Ok(axum::serve(listener, app).await?)
}
