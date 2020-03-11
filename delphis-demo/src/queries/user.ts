import { gql } from "apollo-boost";
import { userProfileFragments } from "./userProfile"

export const userFragments = {
    meInfo: gql`fragment meInfo on User{
        id
        profile {
            ...userProfileFragments.userProfileInfo,
        }
    }`
}

export default {
    me: gql`query GetMe{
        me {
            ...fragments.meInfo,
        }
    }`
};