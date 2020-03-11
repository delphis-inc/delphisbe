import { gql } from "apollo-boost";

export const userProfileFragments = {
    userProfileInfo: gql`fragment profileInfo on UserProfile{
        id
        displayName
        twitterURL {
            displayText
            url
        }
    }`
};