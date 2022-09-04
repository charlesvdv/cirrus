use std::str::FromStr;

use serde::Deserialize;

const ADMIN: &str = "Administrator";
const USER: &str = "User";

#[derive(Clone, Deserialize, Debug)]
pub enum Role {
    Administrator,
    User,
}

impl ToString for Role {
    fn to_string(&self) -> String {
        match self {
            Role::Administrator => ADMIN.to_string(),
            Role::User => USER.to_string(),
        }
    }
}

impl FromStr for Role {
    type Err = std::convert::Infallible;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            ADMIN => Ok(Role::Administrator),
            USER => Ok(Role::User),
            unknown_role => {
                log::error!(
                    "Failed to convert role {:?} to a known role, defaulting to default role",
                    unknown_role
                );
                Ok(Role::User)
            }
        }
    }
}
