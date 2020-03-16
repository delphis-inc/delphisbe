import React, { useState } from 'react';
import classnames from 'classnames';
import './style.scss';
import { useMutation } from '@apollo/react-hooks';
import PostMutations from '../../../mutations/post'

export interface Props {
    discussionID: String | undefined
}

function DiscussionInput(props: Props) {
    const [content, setContent] = useState('');
    const [addPost, { loading: mutationLoading }] = useMutation(PostMutations.createPost);

    const submitPost = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault();
        addPost({ 
            variables: { 
                discussionID: props.discussionID, 
                postContent: content 
            },
        });

        setContent('');
    }

    const buttonTitle = () => {
        if (mutationLoading) {
            return 'Processing...'
        }

        return 'Post';
    }

    return (
        <div className="discussion-input">
            Add to the discussion:
            <div className="flex items-center justify-between">
                <textarea placeholder="Your post here" 
                  className="mr-6 focus:outline-none focus:shadow-outline border rounded-lg py-2 px-4 block w-full appearance-none leading-normal" 
                  name="postContent" 
                  rows={4}
                  value={content} 
                  onChange={(event) => setContent(event.target.value)}
                />
                <button 
                  className={classnames("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded", {'opacity-50 cursor-not-allowed': content === ''})}
                  onClick={submitPost}
                >{buttonTitle()}</button>
            </div>
        </div>
    )
}

export default DiscussionInput;