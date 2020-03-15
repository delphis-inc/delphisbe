import { gql } from 'apollo-boost';

const DiscussionFragments = {
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

export default DiscussionFragments