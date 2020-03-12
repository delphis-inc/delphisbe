import { gql, DocumentNode } from "apollo-boost";

const Fragments = {
    discussionInfo: gql`
    fragment discussionInfo on Discussion{
        id
        anonymityType
        moderator {
            id
            userProfile {
                id
                displayName
                profileImageURL
                twitterURL {
                    displayText
                    url
                }
            }
        }
        participants {
            participantID
        }
        posts {
            id
            isDeleted
            deletedReasonCode
            content
        }
        title
    }`
}

const discussionQuery: DocumentNode = gql`
    query GetDiscussionByID($discussionID: ID!) {
        discussion(id: $discussionID) {
            ...discussionInfo
        }
    }
    ${Fragments.discussionInfo}`

const discussionListQuery: DocumentNode = gql`
    query ListDiscussions {
        listDiscussions {
            ...discussionInfo
        }
    }
    ${Fragments.discussionInfo}`

export { discussionQuery, discussionListQuery }