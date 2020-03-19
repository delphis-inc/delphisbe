import { gql } from "apollo-boost";

export const userProfileFragments = {
    userProfileInfo: gql`fragment userProfileInfo on UserProfile{
        id
        displayName
        twitterURL {
            displayText
            url
        }
    }`
};

export const userFragments = {
    meInfo: gql`fragment meInfo on User{
        id
        profile {
            ...userProfileInfo
        }
    }
    ${userProfileFragments.userProfileInfo}`
}

export default {
    me: gql`query GetMe{
        me {
            ...meInfo
        }
    }
    ${userFragments.meInfo}
    ${userProfileFragments.userProfileInfo}`
};