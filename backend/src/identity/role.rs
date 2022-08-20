use std::str::FromStr;

use serde::Deserialize;

const ADMIN: &'static str = "admin";
const USER: &'static str = "user";

#[derive(Clone, Deserialize, Debug)]
pub enum Role {
    Admin,
    User,
}

impl ToString for Role {
    fn to_string(&self) -> String {
        match self {
            Role::Admin => ADMIN.to_string(),
            Role::User => USER.to_string(),
        }
    }
}

impl FromStr for Role {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            ADMIN => Ok(Role::Admin),
            USER => Ok(Role::User),
            unknown_role => anyhow::bail!("Unknown role '{}'", unknown_role),
        }
    }
}
