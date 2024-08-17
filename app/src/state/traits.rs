pub trait TryFromEnv: Sized {
    fn try_from_env() -> anyhow::Result<Self>;
}
